package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// genDPID generates a distributed process ID (DPID)
// based on a Unix timestamp
func genDPID() string {
	base := "kubed-sh"
	now := time.Now()
	return fmt.Sprintf("%s-%v", base, now.UnixNano())
}

// verify checks if a file exists and returns
// its absolute filepath
func verify(file string) (string, error) {
	fileloc, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(fileloc)
	if err != nil {
		return "", err
	}
	return fileloc, nil
}

// launchhost launches a host, using a template
func launchhost(template, name, image string) error {
	manifestloc := "/tmp/" + name + ".yaml"
	manifest := strings.Replace(template, "DPROC_NAME", name, -1)
	manifest = strings.Replace(manifest, "DPROC_IMAGE", image, -1)
	err := ioutil.WriteFile(manifestloc, []byte(manifest), 0644)
	if err != nil {
		return err
	}
	res, err := kubectl(false, "apply", "-f", manifestloc)
	if err != nil {
		return err
	}
	debug(res)
	return nil
}

// inject uploads the program (binary or script and its dependencies)
// into a given pod and launches it there; for long-running dprocs
// it also creates the service so that we can talk to it via HTTP.
func inject(dproct DProcType, dpid, program, programtype, interpreter, pod string) (string, error) {
	svcname := "undefined"
	// upload program into pod and remember it via annotation:
	src, err := filepath.Abs(program)
	if err != nil {
		return svcname, err
	}
	dest := fmt.Sprintf("%s:/tmp/", pod)
	_, err = kubectl(true, "cp", src, dest)
	if err != nil {
		return svcname, err
	}
	_, err = kubectl(true, "annotate", "pods", pod,
		"original="+program, "interpreter="+interpreter)
	if err != nil {
		return svcname, err
	}
	debug(fmt.Sprintf("Uploaded %s to %s and wrote it into annotation\n", program, pod))
	// start program in pod, handle dproc type specific things:
	var executor string
	execremotefile := fmt.Sprintf("/tmp/%s", program)
	switch interpreter {
	case "binary":
		executor = ""
	default:
		executor = interpreter
	}
	switch dproct {
	case DProcLongRunning:
		// create service for deployment:
		pb := filepath.Base(program)
		svcname = pb[0 : len(pb)-len(filepath.Ext(pb))]
		userdefsvcname := currentenv().evt.get("SERVICE_NAME")
		if userdefsvcname != "" {
			svcname = userdefsvcname
		}
		port := currentenv().evt.get("SERVICE_PORT")
		res, err := kubectl(true, "expose", "deployment", dpid,
			"--name="+svcname, "--port="+port, "--target-port="+port)
		if err != nil {
			return svcname, err
		}
		debug(res)
		// launch program in the given pod, in the background:
		go func() {
			res, err := kubectl(true, "exec", pod, "--", executor, execremotefile)
			if err != nil {
				warn(("Can't start program: " + err.Error()))
			}
			debug(res)
		}()
	case DProcTerminating: // launch program in the given pod
		res, err := kubectl(true, "exec", pod, "--", executor, execremotefile)
		if err != nil {
			warn(("Can't start program: " + err.Error()))
		}
		output(res)
	default:
		return svcname, fmt.Errorf("Can't inject program: unknown distributed process type")
	}
	return svcname, err
}

// launchenv takes a command line input and launches the program
// depending on the dproc type (we consider lines ending with & as long-running)
func launchenv(line, image, interpreter string) (string, string, error) {
	var program, programtype string
	switch interpreter {
	case "binary": // line akin to 'binary args'
		program = strings.Split(line, " ")[0]
		programtype = "bin"
	default: // line akin to 'interpreter script.ext args'
		program = strings.Split(line, " ")[1]
		programtype = "script"
	}
	dpid := genDPID()
	info(fmt.Sprintf("Launching %s of type [%s] with DPID %s", program, programtype, dpid))
	// Step 1. find and verify script locally:
	programloc, err := verify(program)
	if err != nil {
		return dpid, "", err
	}
	_, programfile := filepath.Split(programloc)
	// Step 2. launch host depending on dproc type.
	// If line ends in ' &' we create a Kubernetes
	// deployment and a service, otherwise a pod:
	dproctype := DProcTerminating
	if strings.HasSuffix(line, "&") {
		dproctype = DProcLongRunning
	}
	switch dproctype {
	case DProcLongRunning:
		err := launchhost(longRunningTemplate, dpid, image)
		if err != nil {
			warn("Can't launch long-running distributed process: " + err.Error())
		}
		_, err = kubectl(false, "label", "deployment", dpid,
			"gen=kubed-sh",
			programtype+"="+programfile,
			"env="+currentenv().name,
			"dproctype="+string(dproctype))
		if err != nil {
			warn("Can't label long-running distributed process: " + err.Error())
		}
	case DProcTerminating:
		err := launchhost(terminatingTemplate, dpid, image)
		if err != nil {
			warn("Can't launch terminating distributed process: " + err.Error())
		}
		_, err = kubectl(false, "label", "pod", dpid,
			"gen=kubed-sh",
			programtype+"="+programfile,
			"env="+currentenv().name,
			"dproctype="+string(dproctype))
		if err != nil {
			warn("Can't label terminating distributed process: " + err.Error())
		}
	default:
		warn("Can't launch distributed process, unknown type!")
	}
	time.Sleep(2 * time.Second) // this is a (necessary) hack
	// Step 3. determine the target pod we use to inject program:
	var pod string
	switch dproctype {
	case DProcLongRunning:
		candidatepods, err := kubectl(true, "get", "pods",
			"--selector=app="+dpid, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			debug(err.Error())
		}
		for _, targetpod := range strings.Split(candidatepods, "\n") {
			if strings.HasPrefix(targetpod, dpid) {
				pod = targetpod
				break
			}
		}
	case DProcTerminating:
		pod = dpid
	}
	// Step4. inject program
	svcname, err := inject(dproctype, dpid, program, programtype, interpreter, pod)
	if err != nil {
		warn(fmt.Sprintf("Can't upload and start program: %v", err))
	}
	return dpid, svcname, err
}

func launchbin(line string) (string, string, error) {
	img := currentenv().evt.get("BINARY_IMAGE")
	return launchenv(line, img, "binary")
}

func launchpy(line string) (string, string, error) {
	img := currentenv().evt.get("PYTHON_IMAGE")
	return launchenv(line, img, "python")
}

func launchjs(line string) (string, string, error) {
	img := currentenv().evt.get("NODE_IMAGE")
	return launchenv(line, img, "node")
}

func launchrb(line string) (string, string, error) {
	img := currentenv().evt.get("RUBY_IMAGE")
	return launchenv(line, img, "ruby")
}

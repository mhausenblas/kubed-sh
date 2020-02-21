package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// the time in seconds that the infra process should running
// this process keeps the pod alive, currently 100,000s (~27h)
const keepAliveInSec = "100000"

func gendname() string {
	base := "kubed-sh"
	now := time.Now()
	return fmt.Sprintf("%s-%v", base, now.UnixNano())
}

func extractsrc(line string) string {
	debug("input line: " + line)
	line = strings.TrimSuffix(line, "&")
	line = strings.TrimSpace(line)
	debug("sanitized line: " + line)
	// a binary is standalone:
	if !strings.ContainsAny(line, " ") {
		_, binfile := filepath.Split(line)
		debug("binfile extracted: " + binfile)
		return binfile
	}
	// … otherwise it's a script:
	script := strings.Split(line, " ")[1]
	_, scriptfile := filepath.Split(script)
	debug("scriptfile extracted: " + scriptfile)
	return scriptfile
}

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

func launchenv(line, image, interpreter string) (string, string, error) {
	var binorscript string
	launchtype := "script"
	switch interpreter {
	case "binary":
		// line is something like 'binary args'
		binorscript = strings.Split(line, " ")[0]
		launchtype = "bin"
	default:
		// line is something like 'interpreter script.ext args'
		binorscript = strings.Split(line, " ")[1]
	}
	envdname := gendname()
	deployment := ""
	svcname := ""
	// Step 1. find and verify script locally:
	binorscriptloc, err := verify(binorscript)
	if err != nil {
		return envdname, "", err
	}
	_, binorscriptfile := filepath.Split(binorscriptloc)
	// Step 2. launch interpreter pod:
	// If line ends in a ' &' we create a background
	// distributed process via a deployment and a service,
	// otherwise a simple pod, representing a foreground
	// distributed process:
	dproctype := DProcTerminating
	if strings.HasSuffix(line, "&") {
		dproctype = DProcLongRunning
	}
	go func() {
		res, lerr := kubectl(true, "create", "deployment",
			envdname, "--image="+image, "--", "sh", "-c", "sleep "+keepAliveInSec)
		if lerr != nil {
			warn("Can't launch distributed process: " + lerr.Error())
		}
		info(res)
		_, lerr = kubectl(false, "label", "deployment", envdname,
			"gen=kubed-sh",
			launchtype+"="+binorscriptfile,
			"env="+currentenv().name,
			"dproctype="+string(dproctype))
		if lerr != nil {
			warn("Can't label distributed process: " + lerr.Error())
		}
	}()
	time.Sleep(2 * time.Second) // this is a (necessary) hack
	// set up service and hostpod, if necessary:
	if strings.HasSuffix(line, "&") {
		deployment, err = kubectl(true, "get", "deployment", "--selector=gen=kubed-sh,"+launchtype+"="+binorscriptfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return envdname, "", err
		}
		svcname = binorscriptfile[0 : len(binorscriptfile)-len(filepath.Ext(binorscriptfile))]
		userdefsvcname := currentenv().evt.get("SERVICE_NAME")
		if userdefsvcname != "" {
			svcname = userdefsvcname
		}
		port := currentenv().evt.get("SERVICE_PORT")
		sres, serr := kubectl(true, "expose", "deployment", deployment, "--name="+svcname, "--port="+port, "--target-port="+port)
		if serr != nil {
			return envdname, "", serr
		}
		info(sres)
		candidatepods, serr := kubectl(true, "get", "pods", "--selector=gen=kubed-sh,"+launchtype+"="+binorscriptfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return envdname, "", serr
		}
		for _, canp := range strings.Split(candidatepods, "\n") {
			if strings.HasPrefix(canp, deployment) {
				envdname = canp
				break
			}
		}
	}
	// Step 3. copy script or binary from step 1 into pod and annotate it:
	dest := fmt.Sprintf("%s:/tmp/", envdname)
	_, err = kubectl(true, "cp", binorscriptloc, dest)
	if err != nil {
		return envdname, "", err
	}
	_, err = kubectl(true, "annotate", "pods", envdname, "original="+binorscriptloc, "interpreter="+interpreter)
	if err != nil {
		return envdname, "", err
	}
	info(fmt.Sprintf("uploaded %s to %s\n", binorscriptloc, envdname))
	switch {
	case strings.HasSuffix(line, "&"):
		go func() {
			// Step 4. launch script or binary in pod:
			execremotescript := fmt.Sprintf("/tmp/%s", binorscriptfile)
			err = kubectlbg("exec", "-i", "-t", envdname, interpreter, execremotescript)
			if err != nil {
				debug(err.Error())
			}
		}()
		return deployment, svcname, nil
	default:
		// Step 4. launch script or binary in pod:
		var execres string
		execremotefile := fmt.Sprintf("/tmp/%s", binorscriptfile)
		switch interpreter {
		case "binary":
			execres, err = kubectl(true, "exec", envdname, "--", "sh", "-c", execremotefile)
			if err != nil {
				return envdname, "", err
			}
		default:
			execres, err = kubectl(true, "exec", envdname, interpreter, execremotefile)
			if err != nil {
				return envdname, "", err
			}
		}
		output(execres)
		// Step 5. clean up:
		res, err := kubectl(true, "delete", "pod", envdname)
		if err != nil {
			return envdname, "", err
		}
		debug("delete result " + res)
	}
	return envdname, "", nil
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

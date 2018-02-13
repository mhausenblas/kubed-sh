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

func genpodname() string {
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

func launch(line string) (string, string, error) {
	bin := strings.Split(line, " ")[0]
	hostpod := genpodname()
	var deployment, svcname string
	// Step 1. find and verify binary locally:
	binloc, err := verify(bin)
	if err != nil {
		return hostpod, "", err
	}
	_, binfile := filepath.Split(binloc)
	// Step 2. launch generic pod:
	// If line ends in a ' &' we create a background
	// distributed process via a deployment and a service,
	// otherwise a simple pod, representing a foreground
	// distributed process:
	strategy := "Never"
	dproctype := DProcTerminating
	if strings.HasSuffix(line, "&") {
		strategy = "Always"
		dproctype = DProcLongRunning
	}
	img := currentenv().evt.get("BINARY_IMAGE")
	res, err := kubectl(true, "run", hostpod,
		"--image="+img, "--restart="+strategy,
		"--labels=gen=kubed-sh,bin="+binfile+
			",env="+currentenv().name+
			",dproctype="+string(dproctype),
		"--", "sh", "-c", "sleep "+keepAliveInSec)
	if err != nil {
		return hostpod, "", err
	}
	info(res)
	time.Sleep(5 * time.Second) // this is a hack. need to do prefilght checks and warmup
	// set up service and hostpod, if necessary:
	if strings.HasSuffix(line, "&") {
		deployment, err = kubectl(true, "get", "deployment", "--selector=gen=kubed-sh,bin="+binfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return hostpod, "", err
		}
		svcname = binfile[0 : len(binfile)-len(filepath.Ext(binfile))]
		userdefsvcname := currentenv().evt.get("SERVICE_NAME")
		if userdefsvcname != "" {
			svcname = userdefsvcname
		}
		port := currentenv().evt.get("SERVICE_PORT")
		sres, serr := kubectl(true, "expose", "deployment", deployment, "--name="+svcname, "--port="+port, "--target-port="+port)
		if serr != nil {
			return hostpod, "", serr
		}
		info(sres)
		hostpod, serr = kubectl(true, "get", "pods", "--selector=gen=kubed-sh,bin="+binfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return hostpod, "", serr
		}
	}
	// Step 3. copy binary from step 1 into pod and annotate it:
	dest := fmt.Sprintf("%s:/tmp/", hostpod)
	_, err = kubectl(true, "cp", binloc, dest)
	if err != nil {
		return hostpod, "", err
	}
	_, err = kubectl(true, "annotate", "pods", hostpod, "original="+binloc)
	if err != nil {
		return hostpod, "", err
	}
	info(fmt.Sprintf("Uploaded %s to %s\n", binloc, hostpod))
	// Step 4. launch binary in pod:
	switch {
	case strings.HasSuffix(line, "&"):
		go func() error {
			// Step 4. launch script in pod:
			execremotebin := fmt.Sprintf("/tmp/%s", binfile)
			res, err = kubectl(true, "exec", hostpod, "--", "sh", "-c", execremotebin)
			if err != nil {
				return err
			}
			debug("Exec result " + res)
			return nil
		}()
		return deployment, svcname, nil
	default:
		// Step 4. launch script in pod:
		execremotebin := fmt.Sprintf("/tmp/%s", binfile)
		res, err = kubectl(true, "exec", hostpod, "--", "sh", "-c", execremotebin)
		if err != nil {
			return hostpod, "", err
		}
		output(res)
		// Step 5. clean up:
		_, err = kubectl(true, "delete", "pod", hostpod)
		if err != nil {
			return hostpod, "", err
		}
	}
	return hostpod, "", nil
}

func launchenv(line, image, interpreter string) (string, string, error) {
	// line is something like 'interpreter script.ext args'
	script := strings.Split(line, " ")[1]
	hostpod := genpodname()
	deployment := ""
	svcname := ""
	// Step 1. find and verify script locally:
	scriptloc, err := verify(script)
	if err != nil {
		return hostpod, "", err
	}
	_, scriptfile := filepath.Split(scriptloc)
	// Step 2. launch interpreter pod:
	// If line ends in a ' &' we create a background
	// distributed process via a deployment and a service,
	// otherwise a simple pod, representing a foreground
	// distributed process:
	strategy := "Never"
	dproctype := DProcTerminating
	if strings.HasSuffix(line, "&") {
		strategy = "Always"
		dproctype = DProcLongRunning
	}
	res, err := kubectl(true, "run", hostpod,
		"--image="+image, "--restart="+strategy,
		"--labels=gen=kubed-sh,script="+scriptfile+
			",env="+currentenv().name+
			",dproctype="+string(dproctype),
		"--", "sh", "-c", "sleep "+keepAliveInSec)
	if err != nil {
		return hostpod, "", err
	}
	info(res)
	// this is a hack. need to do prefilght checks and warmup:
	time.Sleep(5 * time.Second)
	// set up service and hostpod, if necessary:
	if strings.HasSuffix(line, "&") {
		deployment, err = kubectl(true, "get", "deployment", "--selector=gen=kubed-sh,script="+scriptfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return hostpod, "", err
		}
		svcname = scriptfile[0 : len(scriptfile)-len(filepath.Ext(scriptfile))]
		userdefsvcname := currentenv().evt.get("SERVICE_NAME")
		if userdefsvcname != "" {
			svcname = userdefsvcname
		}
		port := currentenv().evt.get("SERVICE_PORT")
		sres, serr := kubectl(true, "expose", "deployment", deployment, "--name="+svcname, "--port="+port, "--target-port="+port)
		if serr != nil {
			return hostpod, "", serr
		}
		info(sres)
		hostpod, serr = kubectl(true, "get", "pods", "--selector=gen=kubed-sh,script="+scriptfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return hostpod, "", serr
		}
	}
	// Step 3. copy script from step 1 into pod and annotate it:
	dest := fmt.Sprintf("%s:/tmp/", hostpod)
	_, err = kubectl(true, "cp", scriptloc, dest)
	if err != nil {
		return hostpod, "", err
	}
	_, err = kubectl(true, "annotate", "pods", hostpod, "original="+scriptloc, "interpreter="+interpreter)
	if err != nil {
		return hostpod, "", err
	}
	info(fmt.Sprintf("uploaded %s to %s\n", scriptloc, hostpod))
	switch {
	case strings.HasSuffix(line, "&"):
		go func() error {
			// Step 4. launch script in pod:
			execremotescript := fmt.Sprintf("/tmp/%s", scriptfile)
			err = kubectlbg("exec", "-i", "-t", hostpod, interpreter, execremotescript)
			if err != nil {
				return err
			}
			return nil
		}()
		return deployment, svcname, nil
	default:
		// Step 4. launch script in pod:
		execremotescript := fmt.Sprintf("/tmp/%s", scriptfile)
		res, err = kubectl(true, "exec", hostpod, interpreter, execremotescript)
		if err != nil {
			return hostpod, "", err
		}
		output(res)
		// Step 5. clean up:
		_, err = kubectl(true, "delete", "pod", hostpod)
		if err != nil {
			return hostpod, "", err
		}
	}
	return hostpod, "", nil
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

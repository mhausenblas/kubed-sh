package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func genpodname() string {
	base := "kubed-sh"
	now := time.Now()
	return fmt.Sprintf("%s-%v", base, now.UnixNano())
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

func launch(binary string) error {
	hostpod := genpodname()
	// Step 1. find and verify binary locally:
	binloc, err := verify(binary)
	if err != nil {
		return err
	}
	// Step 2. launch generic pod:
	res, err := kubectl("run", hostpod, "--image=alpine:3.7", "--restart=Never", "--", "sh", "-c", "sleep 10000")
	if err != nil {
		return err
	}
	info(res)
	time.Sleep(2 * time.Second) // this is a hack. need to do prefilght checks and warmup
	// Step 3. copy binary from step 1 into pod:
	dest := fmt.Sprintf("%s:/tmp/", hostpod)
	_, err = kubectl("cp", binloc, dest)
	if err != nil {
		return err
	}
	info(fmt.Sprintf("Uploaded %s to %s\n", binloc, hostpod))
	// Step 4. launch binary in pod:
	_, binfile := filepath.Split(binloc)
	execremotebin := fmt.Sprintf("/tmp/%s", binfile)
	res, err = kubectl("exec", hostpod, "--", "sh", "-c", execremotebin)
	if err != nil {
		return err
	}
	output(res)
	// Step 5. clean up:
	_, err = kubectl("delete", "pod", hostpod)
	if err != nil {
		return err
	}
	return nil
}

func launchenv(line, image, interpreter string) error {
	// line is something like 'interpreter script.ext args'
	script := strings.Split(line, " ")[1]
	hostpod := genpodname()
	// Step 1. find and verify script locally:
	scriptloc, err := verify(script)
	if err != nil {
		return err
	}
	_, scriptfile := filepath.Split(scriptloc)
	// Step 2. launch interpreter pod:
	// If line ends in a ' &' we create a background
	// distributed process via a deployment and a service,
	// otherwise a simple pod, representing a foreground
	// distributed process:
	strategy := "Never"
	if strings.HasSuffix(line, "&") {
		strategy = "Always"
	}
	res, err := kubectl("run", hostpod,
		"--image="+image, "--restart="+strategy,
		"--labels=gen=kubed-sh,script="+scriptfile,
		"--", "sh", "-c", "sleep 10000")
	if err != nil {
		return err
	}
	info(res)
	// this is a hack. need to do prefilght checks and warmup:
	time.Sleep(2 * time.Second)
	// set up service and hostpod, if necessary:
	if strings.HasSuffix(line, "&") {
		deployment, serr := kubectl("get", "deployment", "--selector=gen=kubed-sh,script="+scriptfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return serr
		}
		svcname := scriptfile[0 : len(scriptfile)-len(filepath.Ext(scriptfile))]
		sres, serr := kubectl("expose", "deployment", deployment, "--name="+svcname, "--port=80", "--target-port=80")
		if serr != nil {
			return serr
		}
		info(sres)
		hostpod, serr = kubectl("get", "pods", "--selector=gen=kubed-sh,script="+scriptfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return serr
		}
	}
	// Step 3. copy script from step 1 into pod:
	dest := fmt.Sprintf("%s:/tmp/", hostpod)
	_, err = kubectl("cp", scriptloc, dest)
	if err != nil {
		return err
	}
	info(fmt.Sprintf("Uploaded %s to %s\n", scriptloc, hostpod))
	// Step 4. launch script in pod:
	execremotescript := fmt.Sprintf("/tmp/%s", scriptfile)
	res, err = kubectl("exec", hostpod, interpreter, execremotescript)
	if err != nil {
		return err
	}
	output(res)
	// Step 5. clean up:
	_, err = kubectl("delete", "pod", hostpod)
	if err != nil {
		return err
	}
	return nil
}

func launchpy(line string) error {
	return launchenv(line, "python:3.6-alpine3.7", "python")
}

func launchjs(line string) error {
	return launchenv(line, "node:9.4-alpine", "node")
}

func launchrb(line string) error {
	return launchenv(line, "ruby:2.5-alpine3.7", "ruby")
}

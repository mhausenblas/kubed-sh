package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func genpodname() string {
	base := "kubed-sh"
	now := time.Now()
	return fmt.Sprintf("%s-%v", base, now.UnixNano())
}

func launch(binary string) error {
	hostpod := genpodname()
	// Step 1. find and verify binary locally:
	binloc, err := filepath.Abs(binary)
	if err != nil {
		return err
	}
	_, err = os.Stat(binloc)
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

func launchpy(script string) error {
	hostpod := genpodname()
	// Step 1. find and verify Python script locally:
	scriptloc, err := filepath.Abs(script)
	if err != nil {
		return err
	}
	_, err = os.Stat(scriptloc)
	if err != nil {
		return err
	}
	// Step 2. launch Python pod:
	res, err := kubectl("run", hostpod, "--image=python:3.6-alpine3.7", "--restart=Never", "--", "sh", "-c", "sleep 10000")
	if err != nil {
		return err
	}
	info(res)
	time.Sleep(2 * time.Second) // this is a hack. need to do prefilght checks and warmup
	// Step 3. copy script from step 1 into pod:
	dest := fmt.Sprintf("%s:/tmp/", hostpod)
	_, err = kubectl("cp", scriptloc, dest)
	if err != nil {
		return err
	}
	info(fmt.Sprintf("Uploaded %s to %s\n", scriptloc, hostpod))
	// Step 4. launch script in pod:
	_, scriptfile := filepath.Split(scriptloc)
	execremotescript := fmt.Sprintf("/tmp/%s", scriptfile)
	res, err = kubectl("exec", hostpod, "python", execremotescript)
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

func launchjs(script string) error {
	hostpod := genpodname()
	// Step 1. find and verify Node.js script locally:
	scriptloc, err := filepath.Abs(script)
	if err != nil {
		return err
	}
	_, err = os.Stat(scriptloc)
	if err != nil {
		return err
	}
	// Step 2. launch Node.js pod:
	res, err := kubectl("run", hostpod, "--image=node:9.4-alpine", "--restart=Never", "--", "sh", "-c", "sleep 10000")
	if err != nil {
		return err
	}
	info(res)
	time.Sleep(2 * time.Second) // this is a hack. need to do prefilght checks and warmup
	// Step 3. copy script from step 1 into pod:
	dest := fmt.Sprintf("%s:/tmp/", hostpod)
	_, err = kubectl("cp", scriptloc, dest)
	if err != nil {
		return err
	}
	info(fmt.Sprintf("Uploaded %s to %s\n", scriptloc, hostpod))
	// Step 4. launch script in pod:
	_, scriptfile := filepath.Split(scriptloc)
	execremotescript := fmt.Sprintf("/tmp/%s", scriptfile)
	res, err = kubectl("exec", hostpod, "node", execremotescript)
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

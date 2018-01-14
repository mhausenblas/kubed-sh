package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// output prints primary, output infor to shell
func output(msg string) {
	fmt.Println(msg)
}

// info prints secondary, non-output info to shell
func info(msg string) {
	fmt.Printf("\033[34m%s\033[0m\n", msg)
}

func checkruntime() {
	switch runtime.GOOS {
	case "linux":
		fmt.Printf("Note: As you're running kubed-sh on Linux you can directly launch binaries.\n\n")
	default:
		fmt.Printf("Note: It seems you're running kubed-sh in a non-Linux environment (detected: %s),\nso make sure the binaries you launch are Linux binaries in ELF format.\n\n", runtime.GOOS)
	}
}

func shellout(cmd string, args ...string) (string, error) {
	result := ""
	var out bytes.Buffer
	log.Debug(cmd, args)
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		return result, err
	}
	result = strings.TrimSpace(out.String())
	return result, nil
}

func kubectl(cmd string, args ...string) (string, error) {
	kubectlbin, err := shellout("which", "kubectl")
	if err != nil {
		return "", err
	}
	all := append([]string{cmd}, args...)
	result, err := shellout(kubectlbin, all...)
	if err != nil {
		return "", err
	}
	return result, nil
}

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
	res, err := kubectl("run", hostpod, "--image=python:3.6", "--restart=Never", "--", "sh", "-c", "sleep 10000")
	if err != nil {
		return err
	}
	info(res)
	time.Sleep(2 * time.Second) // this is a hack. need to do prefilght checks and warmup
	// Step 3. copy binary from step 1 into pod:
	dest := fmt.Sprintf("%s:/tmp/", hostpod)
	_, err = kubectl("cp", scriptloc, dest)
	if err != nil {
		return err
	}
	info(fmt.Sprintf("Uploaded %s to %s\n", scriptloc, hostpod))
	// Step 4. launch binary in pod:
	_, binfile := filepath.Split(scriptloc)
	execremotescript := fmt.Sprintf("/tmp/%s", binfile)
	res, err = kubectl("exec", hostpod, "--", "python", execremotescript)
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

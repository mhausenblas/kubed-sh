package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// output prints primary, output messages to shell
func output(msg string) {
	fmt.Println(msg)
}

// info prints secondary, non-output info to shell
func info(msg string) {
	fmt.Printf("\033[34m%s\033[0m\n", msg)
}

// warn prints warning messages to shell
func warn(msg string) {
	fmt.Printf("\033[31m%s\033[0m\n", msg)
}

// debug prints debug messages to shell
func debug(msg string) {
	if DEBUG {
		fmt.Printf("\033[33m%s\033[0m\n", msg)
	}
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
	if DEBUG {
		c.Stderr = os.Stderr
	}
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		return result, err
	}
	result = strings.TrimSpace(out.String())
	return result, nil
}

func shelloutbg(cmd string, args ...string) error {
	log.Debug(cmd, args)
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
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

func kubectlbg(cmd string, args ...string) error {
	kubectlbin, err := shellout("which", "kubectl")
	if err != nil {
		return err
	}
	all := append([]string{cmd}, args...)
	err = shelloutbg(kubectlbin, all...)
	if err != nil {
		return err
	}
	return nil
}

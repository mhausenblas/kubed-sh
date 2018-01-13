package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

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

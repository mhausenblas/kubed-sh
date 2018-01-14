package main

import (
	"fmt"
	"strconv"
	"strings"
)

func launchfail(line, reason string) {
	fmt.Printf("\nFailed to launch %s in the cluster due to:\n%s\n\n", strconv.Quote(line), reason)
	husage(line)
}

func hlaunch(line string) {
	// If a line doesn't start with one of the
	// known environments, assume user wants to
	// launch a binary:
	switch {
	case strings.HasPrefix(line, "python "):
		err := launchpy(line)
		if err != nil {
			launchfail(line, err.Error())
		}
	case strings.HasPrefix(line, "node "):
		err := launchjs(line)
		if err != nil {
			launchfail(line, err.Error())
		}
	case strings.HasPrefix(line, "ruby "):
		err := launchrb(line)
		if err != nil {
			launchfail(line, err.Error())
		}
	default:
		err := launch(line)
		if err != nil {
			launchfail(line, err.Error())
		}
	}
}

func hliterally(line string) {
	l := strings.Split(line, " ")
	res, err := kubectl(l[1], l[2:]...)
	if err != nil {
		fmt.Printf("\nFailed to execute kubectl %s command due to:\n%s\n\n", l[1:], err)
	}
	output(res)
}

func hecho(line string) {
	l := strings.Split(line, " ")
	fmt.Println(l[1])
}

func husage(line string) {
	fmt.Println("The available built-in commands of kubed-sh are:")
	fmt.Printf("%s", completer.Tree("    "))
	fmt.Println("\nTo run a program in the Kubernetes cluster, simply specify the binary\nor call it with one of the following supported interpreters:")
	fmt.Printf("    - Node.js … node script.js (default version: 9.4)\n    - Python … python script.py (default version: 3.6)\n    - Ruby … ruby script.rb (default version: 2.5)\n")
}

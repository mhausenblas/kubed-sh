package main

import (
	"fmt"
	"strconv"
	"strings"
)

func hlaunch(line string) {
	// If a line doesn't start with one of the
	// known environments, assume user wants to
	// launch a binary, so try to find and copy this
	// into the pod, then execute it:
	l := strings.Split(line, " ")
	switch {
	case strings.HasPrefix(line, "python "):
		err := launchpy(l[1])
		if err != nil {
			fmt.Printf("\nFailed to launch %s in the cluster due to:\n%s\n\n", strconv.Quote(line), err)
			husage(line)
		}
	case strings.HasPrefix(line, "node "):
		fmt.Printf("Launching a node:9-image based container and executing %s in it\n", l[1])
	case strings.HasPrefix(line, "ruby "):
		fmt.Printf("Launching a ruby:2.5-image based container and executing %s in it\n", l[1])
	default:
		err := launch(line)
		if err != nil {
			fmt.Printf("\nFailed to launch %s in the cluster due to:\n%s\n\n", strconv.Quote(line), err)
			husage(line)
		}
	}
}

func hliterally(line string) {
	l := strings.Split(line, " ")
	_, err := kubectl(l[1], l[1:]...)
	if err != nil {
		fmt.Printf("\nFailed to execute kubectl command due to:\n%s\n\n", strconv.Quote(line), err)
	}
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

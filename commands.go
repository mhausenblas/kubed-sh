package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

func huse(line string, rl *readline.Instance) {
	if !strings.ContainsAny(line, " ") {
		info("Need a target cluster")
		return
	}
	targetcontext := strings.Split(line, " ")[1]
	res, err := kubectl("config", "use-context", targetcontext)
	if err != nil {
		fmt.Printf("\nFailed to switch contexts due to:\n%s\n\n", err)
		return
	}
	output(res)
	rl.SetPrompt(fmt.Sprintf("[\033[32m%s\033[0m]$ ", targetcontext))
}

func hcontexts() {
	res, err := kubectl("config", "get-contexts")
	if err != nil {
		fmt.Printf("\nFailed to list contexts due to:\n%s\n\n", err)
	}
	output(res)
}

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
	if !strings.ContainsAny(line, " ") {
		info("Not enough input for a valid kubectl command")
		return
	}
	l := strings.Split(line, " ")
	res, err := kubectl(l[1], l[2:]...)
	if err != nil {
		fmt.Printf("\nFailed to execute kubectl %s command due to:\n%s\n\n", l[1:], err)
	}
	output(res)
}

func hecho(line string) {
	if !strings.ContainsAny(line, " ") {
		info("No value to echo given")
		return
	}
	l := strings.Split(line, " ")
	fmt.Println(l[1])
}

func husage(line string) {
	fmt.Println("The available built-in commands of kubed-sh are:")
	fmt.Printf("%s", completer.Tree("    "))
	fmt.Println("\nTo run a program in the Kubernetes cluster, simply specify the binary\nor call it with one of the following supported interpreters:")
	fmt.Printf("    - Node.js … node script.js (default version: 9.4)\n    - Python … python script.py (default version: 3.6)\n    - Ruby … ruby script.rb (default version: 2.5)\n")
}

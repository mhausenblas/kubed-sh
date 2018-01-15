package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

func henv() {
	for k, v := range evt.et {
		output(fmt.Sprintf("%s=%s", k, v))
	}
}

func hkill(line string) {
	if !strings.ContainsAny(line, " ") {
		info("Need a target distributed process to kill")
		return
	}
	// pre-flight check if dproc exists:
	ID := strings.Split(line, " ")[1]
	_, err := kubectl("get", "deployment", ID)
	if err != nil {
		warn(fmt.Sprintf("A distributed process with the ID '%s' does not exist in current context. Try the ps command …", ID))
		return
	}
	_, err = kubectl("scale", "--replicas=0", "deployment", ID)
	if err != nil {
		killfail(line, err.Error())
		return
	}
	_, err = kubectl("delete", "deployment", ID)
	if err != nil {
		killfail(line, err.Error())
		return
	}

	// gather info to remove from global DPT:
	kubecontext, err := kubectl("config", "current-context")
	if err != nil {
		killfail(line, err.Error())
		return
	}
	dproc, err := dpt.getDProc(ID, kubecontext)
	if err != nil {
		killfail(line, err.Error())
		return
	}
	// something like xxx:blah
	src := strings.Split(dproc.Src, ":")[1]
	// now get rid of the extension:
	svcname := src[0 : len(src)-len(filepath.Ext(src))]
	_, err = kubectl("delete", "service", svcname)
	if err != nil {
		killfail(line, err.Error())
		return
	}
	if err != nil {
		info(err.Error())
	}
	// finally, remove drpoc from global DPT:
	dpt.removeDProc(dproc)
}

func hps(line string) {
	args := ""
	if strings.ContainsAny(line, " ") {
		args = strings.Split(line, " ")[1]
	}
	var kubecontext string
	switch args {
	case "all":
		kubecontext = ""
	default:
		k, err := kubectl("config", "current-context")
		if err != nil {
			warn("Can't determine current context")
			return
		}
		kubecontext = k
	}
	res := dpt.DumpDPT(kubecontext)
	output(res)
}

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
}

func killfail(line, reason string) {
	fmt.Printf("\nFailed to kill %s due to:\n%s\n\n", strconv.Quote(line), reason)
}

func hlaunch(line string) {
	// If a line doesn't start with one of the
	// known environments, assume user wants to
	// launch a binary:
	var dpid string
	src := extractsrc(line)
	src = "script:" + src
	switch {
	case strings.HasPrefix(line, "python "):
		d, err := launchpy(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
	case strings.HasPrefix(line, "node "):
		d, err := launchjs(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
	case strings.HasPrefix(line, "ruby "):
		d, err := launchrb(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
	default: // binary
		d, err := launch(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
		src = "bin:" + extractsrc(line)
	}
	// update DPT
	if strings.HasSuffix(line, "&") {
		kubecontext, err := kubectl("config", "current-context")
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpt.addDProc(newDProc(dpid, DProcLongRunning, kubecontext, src))
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
	echo := strings.Split(line, " ")[1]
	if strings.HasPrefix(echo, "$") {
		v := evt.get(echo[1:])
		if v != "" {
			fmt.Println(v)
			return
		}
	}
	fmt.Println(echo)
}

func husage(line string) {
	fmt.Println("The available built-in commands of kubed-sh are:")
	fmt.Printf("%s", completer.Tree("    "))
	fmt.Println("\nTo run a program in the Kubernetes cluster, simply specify the binary\nor call it with one of the following supported interpreters:")
	fmt.Printf("    - Node.js … node script.js (default version: 9.4)\n    - Python … python script.py (default version: 3.6)\n    - Ruby … ruby script.rb (default version: 2.5)\n")
}

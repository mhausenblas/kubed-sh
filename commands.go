package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

func henv() {
	tmp := []string{}
	for k, v := range evt.et {
		tmp = append(tmp, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(tmp)
	for _, e := range tmp {
		output(e)
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
	if !strings.ContainsAny(line, " ") {
		helpall()
		return
	}
	cmd := strings.Split(line, " ")[1]
	switch {
	case cmd == "contexts":
		cmd += "\n\nThis is a local command that lists all currently available Kubernetes contexts you can work with.\nA context is a (cluster, namespace, user) tuple, see also https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/"
	case cmd == "echo":
		cmd += " val\n\nThis is a local command that prints the literal value 'val' or an environment variable if prefixed with an '$'."
	case cmd == "env":
		cmd += "\n\nThis is a local command that lists all environment variables currently defined."
	case cmd == "exit" || cmd == "quit":
		cmd += "\n\nThis is a local command that you can use to leave the kubed-sh shell."
	case cmd == "help":
		cmd += "\n\nThis is a local command that lists all built-in commands. You can use 'help command' for more details on a certain command."
	case cmd == "kill":
		cmd += " $DPID\n\nThis is a cluster command that stops the distributed process with the distributed process ID 'DPID'."
	case cmd == "literally":
		cmd += " $COMMAND\n\nThis is a local command that executes what follows as a kubectl command\n.Note that you can also prefix a line with ` to achieve the same."
	case cmd == "ps":
		cmd += " [all]\n\nThis is a cluster command that lists all distributed (long-running) processes in the current context.\nIf used with the optional 'all' argument then all distributed processes across all contexts are shown."
	case cmd == "pwd":
		cmd += "\n\nThis is a local command that prints the current working directory on your local machine."
	case cmd == "use":
		cmd += " $CONTEXT\n\nThis is a local command that selects a certain context to work with. \nA context is a (cluster, namespace, user) tuple, see also https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/"
	default:
		cmd += "\n\nNo details available, yet."
	}
	fmt.Println(cmd)

}

func helpall() {
	fmt.Printf("The available built-in commands of kubed-sh are:\n\n")
	for _, e := range completer.GetChildren() {
		cmd := strings.TrimSpace(string(e.GetName()))
		switch {
		case cmd == "contexts":
			cmd += " (local):\n\t\tlist available Kubernetes contexts (cluster, namespace, user tuples)"
		case cmd == "echo":
			cmd += " (local):\n\t\tprint a value or environment variable"
		case cmd == "env":
			cmd += " (local):\n\t\tlist all environment variables currently defined"
		case cmd == "exit" || cmd == "quit":
			cmd += " (local):\n\t\tleave shell"
		case cmd == "help":
			cmd += " (local):\n\t\tlist built-in commands; use help command for more details"
		case cmd == "kill":
			cmd += " (cluster):\n\t\tstop a distributed process"
		case cmd == "literally":
			cmd += " (local):\n\t\texecute what follows as a kubectl command\n\t\tnote that you can also prefix a line with ` to achieve the same"
		case cmd == "ps":
			cmd += " (cluster):\n\t\tlist all distributed (long-running) processes in current context"
		case cmd == "pwd":
			cmd += " (local):\n\t\tprint current working directory"
		case cmd == "use":
			cmd += " (local):\n\t\tselect a certain context to work with"
		default:
			cmd += "TBD"
		}
		fmt.Println(cmd)
	}
	fmt.Println(`
To run a program in the Kubernetes cluster, specify the binary
or call it with one of the following supported interpreters:

- Node.js … node script.js (default version: 9.4)
- Python … python script.py (default version: 3.6)
- Ruby … ruby script.rb (default version: 2.5)
`)
}

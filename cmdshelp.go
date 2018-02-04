package main

import (
	"fmt"
	"strings"
)

func helpall() {
	fmt.Printf("The available built-in commands of kubed-sh are:\n\n")
	for _, e := range completer.GetChildren() {
		cmd := strings.TrimSpace(string(e.GetName()))
		switch {
		case cmd == "cat":
			cmd += " (local):\n\t\t\toutput content of file to terminal"
		case cmd == "contexts":
			cmd += " (local):\n\t\t\tlist available Kubernetes contexts (cluster, namespace, user tuples)"
		case cmd == "cd":
			cmd += " (local):\n\t\t\tchange working directory"
		case cmd == "curl":
			cmd += " (cluster):\n\t\t\texecute a curl operation in the cluster"
		case cmd == "echo":
			cmd += " (local):\n\t\t\tprint a value or environment variable"
		case cmd == "env":
			cmd += " (local):\n\t\t\tlist all environment variables currently defined"
		case cmd == exitcmd:
			cmd += " (local):\n\t\t\tleave shell"
		case cmd == "help":
			cmd += " (local):\n\t\t\tlist built-in commands; use help command for more details"
		case cmd == "kill":
			cmd += " (cluster):\n\t\t\tstop a distributed process"
		case cmd == "literally":
			cmd += " (local):\n\t\t\texecute what follows as a kubectl command\n\t\t\tnote that you can also prefix a line with ` to achieve the same"
		case cmd == "ls":
			cmd += " (local):\n\t\t\tlists content of directory"
		case cmd == "ps":
			cmd += " (cluster):\n\t\t\tlist all distributed (long-running) processes in current context"
		case cmd == "pwd":
			cmd += " (local):\n\t\t\tprint current working directory"
		case cmd == "sleep":
			cmd += " (local):\n\t\t\tsleep for specified time interval (NOP)"
		case cmd == "use":
			cmd += " (local):\n\t\t\tselect a certain context to work with"
		default:
			cmd += "\t\tto be done"
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

func husage(line string) {
	if !strings.ContainsAny(line, " ") {
		helpall()
		return
	}
	cmd := strings.Split(line, " ")[1]
	switch {
	case cmd == "cat":
		cmd += " $filename\n\nThis is a local command that outputs the content of file 'filename' to the terminal."
	case cmd == "contexts":
		cmd += "\n\nThis is a local command that lists all currently available Kubernetes contexts you can work with.\nA context is a (cluster, namespace, user) tuple, see also https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/"
	case cmd == "cd":
		cmd += " $dir\n\nThis is a local command that changes the current directory to 'dir'."
	case cmd == "curl":
		cmd += " $URL\n\nThis is a cluster command that executes curl against the URL 'URL'."
	case cmd == "echo":
		cmd += " val\n\nThis is a local command that prints the literal value 'val' or an environment variable if prefixed with an '$'."
	case cmd == "env":
		cmd += "\n\nThis is a local command that lists all environment variables currently defined."
	case cmd == exitcmd:
		cmd += "\n\nThis is a local command that you can use to leave the kubed-sh shell."
	case cmd == "help":
		cmd += "\n\nThis is a local command that lists all built-in commands. You can use 'help command' for more details on a certain command."
	case cmd == "kill":
		cmd += " $DPID\n\nThis is a cluster command that stops the distributed process with the distributed process ID 'DPID'."
	case cmd == "literally":
		cmd += " $COMMAND\n\nThis is a local command that executes what follows as a kubectl command\n.Note that you can also prefix a line with ` to achieve the same."
	case cmd == "ls":
		cmd += " $dir\n\nThis is a local command that lists the content of directory 'dir'."
	case cmd == "ps":
		cmd += " [all]\n\nThis is a cluster command that lists all distributed (long-running) processes in the current context.\nIf used with the optional 'all' argument then all distributed processes across all contexts are shown."
	case cmd == "pwd":
		cmd += "\n\nThis is a local command that prints the current working directory on your local machine."
	case cmd == "sleep":
		cmd += " $TIME_INTERVAL\n\nThis is a local command that pauses execution for the specified time interval, for example 'sleep 3s' or 'sleep 450ms'. \nFor formatting, see also https://golang.org/pkg/time/#ParseDuration"
	case cmd == "use":
		cmd += " $CONTEXT\n\nThis is a local command that selects a certain context to work with. \nA context is a (cluster, namespace, user) tuple, see also \nhttps://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/"
	default:
		cmd += "\n\nNo details available, yet."
	}
	fmt.Println(cmd)
}

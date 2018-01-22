package main

import (
	"fmt"
	"strings"
)

func husage(line string) {
	if !strings.ContainsAny(line, " ") {
		helpall()
		return
	}
	cmd := strings.Split(line, " ")[1]
	switch {
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
	case cmd == "exit":
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
		case cmd == "cd":
			cmd += " (local):\n\t\tchanges working directory"
		case cmd == "curl":
			cmd += " (cluster):\n\t\texecutes a curl operation in the cluster"
		case cmd == "echo":
			cmd += " (local):\n\t\tprint a value or environment variable"
		case cmd == "env":
			cmd += " (local):\n\t\tlist all environment variables currently defined"
		case cmd == "exit":
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

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("help"),
	readline.PcItem("ps"),
	readline.PcItem("pwd"),
)

func main() {
	checkruntime()
	kubecontext, err := kubectl("config", "current-context")
	if err != nil {
		panic(err)
	}
	rl, err := readline.NewEx(&readline.Config{
		AutoComplete:    completer,
		Prompt:          fmt.Sprintf("[\033[32m%s\033[0m]$ ", kubecontext),
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = rl.Close()

	}()
	log.SetOutput(rl.Stderr())
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		switch {
		case line == "help":
			husage(line)
		case strings.HasPrefix(line, "echo"):
			hecho(line)
		case strings.HasPrefix(line, "literally") || strings.HasPrefix(line, "`"):
			hliterally(line)
		case strings.HasPrefix(line, "ps"):
			fmt.Println("listing your distributed processes running in the cluster")
		case strings.HasPrefix(line, "pwd"):
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Can't determine where I am due to:\n%s", err)
			}
			fmt.Println(cwd)
		case line == "exit" || line == "quit":
			goto exit
		case line == "":
		default:
			hlaunch(line)
		}
	}
exit:
}

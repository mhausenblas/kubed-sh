package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
)

func main() {
	kubecontext, err := kubectl("config", "current-context")
	if err != nil {
		panic(err)
	}
	rl, err := readline.NewEx(&readline.Config{
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
		case line == "exit" || line == "quit":
			goto exit
		default:
			fmt.Println("unknown command", strconv.Quote(line))
		}
	}
exit:
}

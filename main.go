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

	kubecontext := "replace me with kube context> "

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          kubecontext,
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
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
		case line == "exit":
			goto exit
		default:
			fmt.Println("unknown command", strconv.Quote(line))
		}
	}
exit:
}

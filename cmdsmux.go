package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
)

func interpret(rl *readline.Instance) {
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
		case strings.HasPrefix(line, "contexts"):
			hcontexts()
		case strings.HasPrefix(line, "cd"):
			hcd(line)
		case strings.HasPrefix(line, "curl"):
			hcurl(line)
		case strings.HasPrefix(line, "echo"):
			hecho(line)
		case strings.HasPrefix(line, "env"):
			henv()
		case strings.HasPrefix(line, "help"):
			husage(line)
		case strings.HasPrefix(line, "kill"):
			hkill(line)
		case strings.HasPrefix(line, "literally") || strings.HasPrefix(line, "`"):
			if strings.HasPrefix(line, "`") {
				line = fmt.Sprintf("literally %s", strings.TrimPrefix(line, "`"))
			}
			hliterally(line)
		case strings.HasPrefix(line, "cat"):
			hlocalexec(line)
		case strings.HasPrefix(line, "ls"):
			hlocalexec(line)
		case strings.HasPrefix(line, "ps"):
			hps(line)
		case strings.HasPrefix(line, "pwd"):
			hlocalexec(line)
		case strings.HasPrefix(line, "use"):
			huse(line, rl)
		case line == "debug":
			switch debugmode {
			case true:
				debugmode = false
			case false:
				debugmode = true
			}
		case line == "exit":
			return
		case line == "version":
			output(releaseVersion)
		case strings.Contains(line, "="):
			envar := strings.Split(line, "=")[0]
			value := strings.Split(line, "=")[1]
			evt.set(envar, value)
		case line == "" || strings.HasPrefix(line, "#"):
		default:
			hlaunch(line)
		}
	}
}

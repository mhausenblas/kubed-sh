package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
)

// interpret interactively interprets commands
func interpreti(rl *readline.Instance) {
	for {
		line, err := rl.Readline()
		if err != nil {
			switch err {
			case readline.ErrInterrupt:
				debug("got interrupted")
				continue
			case io.EOF:
				warn("got an EOF")
			default:
				debug("unkown signal received: " + err.Error())
			}
		}
		line = strings.TrimSpace(line)
		done := interpretl(line)
		if done {
			return
		}
	}
}

// interprets interprets a script,
// line by line
func interprets(script string) {
	lines := strings.Split(script, "\n")
	for _, line := range lines {
		done := interpretl(line)
		if done {
			return
		}
	}
}

// interpretl interprets a single line,
// used both in interactive and scripting mode.
func interpretl(line string) bool {
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
		henv(line)
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
	case strings.HasPrefix(line, "sleep"):
		hsleep(line)
	case strings.HasPrefix(line, "use"):
		huse(line)
	case line == "debug":
		switch debugmode {
		case true:
			debugmode = false
			info("DEBUG mode is now off.")
		case false:
			debugmode = true
			info("DEBUG mode is now on.")
		}
	case line == exitcmd:
		return true
	case line == "version":
		output(version)
	case strings.Contains(line, "="):
		envar := strings.Split(line, "=")[0]
		value := strings.Split(line, "=")[1]
		currentenv().evt.set(envar, value)
	case line == "" || strings.HasPrefix(line, "#"):
	default:
		hlaunch(line)
	}
	return false
}

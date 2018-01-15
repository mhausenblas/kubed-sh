package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
)

var (
	debugmode bool
	completer = readline.NewPrefixCompleter(
		readline.PcItem("contexts"),
		readline.PcItem("echo"),
		readline.PcItem("env"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
		readline.PcItem("kill"),
		readline.PcItem("literally"),
		readline.PcItem("ps"),
		readline.PcItem("pwd"),
		readline.PcItem("use"),
		readline.PcItem("quit"),
	)
)

func init() {
	if envd := os.Getenv("DEBUG"); envd != "" {
		debugmode = true
	}
	dpt = &DProcTable{
		mux: new(sync.Mutex),
		lt:  make(map[string]DProc),
	}
	err := dpt.BuildDPT()
	if err != nil {
		output(err.Error())
	}
	evt = &EnvVarTable{
		mux: new(sync.Mutex),
		et:  make(map[string]string),
	}
}

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
		case strings.HasPrefix(line, "contexts"):
			hcontexts()
		case strings.HasPrefix(line, "echo"):
			hecho(line)
		case strings.HasPrefix(line, "env"):
			henv()
		case line == "help":
			husage(line)
		case strings.HasPrefix(line, "kill"):
			hkill(line)
		case strings.HasPrefix(line, "literally") || strings.HasPrefix(line, "`"):
			if strings.HasPrefix(line, "`") {
				line = fmt.Sprintf("literally %s", strings.TrimPrefix(line, "`"))
			}
			hliterally(line)
		case strings.HasPrefix(line, "ps"):
			hps(line)
		case strings.HasPrefix(line, "pwd"):
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Can't determine where I am due to:\n%s", err)
			}
			fmt.Println(cwd)
		case strings.HasPrefix(line, "use"):
			huse(line, rl)
		case line == "exit" || line == "quit":
			goto exit
		case strings.Contains(line, "="):
			envar := strings.Split(line, "=")[0]
			value := strings.Split(line, "=")[1]
			evt.set(envar, value)
		case line == "":
		default:
			hlaunch(line)
		}
	}
exit:
}

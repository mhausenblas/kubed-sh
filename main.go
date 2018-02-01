package main

import (
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
)

var (
	releaseVersion string
	debugmode      bool
	noprepull      bool
	customkubectl  string
	prevdir        string
	rl             *readline.Instance
	completer      = readline.NewPrefixCompleter(
		readline.PcItem("cat"),
		readline.PcItem("cd"),
		readline.PcItem("curl"),
		readline.PcItem("contexts"),
		readline.PcItem("echo"),
		readline.PcItem("env"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
		readline.PcItem("kill"),
		readline.PcItem("literally"),
		readline.PcItem("ls"),
		readline.PcItem("ps"),
		readline.PcItem("pwd"),
		readline.PcItem("use"),
	)
)

func init() {
	if env := os.Getenv("KUBEDSH_DEBUG"); env != "" {
		debugmode = true
	}
	if env := os.Getenv("KUBEDSH_NOPREPULL"); env != "" {
		noprepull = true
	}
	if env := os.Getenv("KUBECTL_BINARY"); env != "" {
		customkubectl = env
	}
	prevdir, _ = os.Getwd()
	// set up the global distributed process table:
	dpt = &DProcTable{
		mux: new(sync.Mutex),
		lt:  make(map[string]DProc),
	}
	err := dpt.BuildDPT()
	if err != nil {
		output(err.Error())
	}
	// set up the environment variables table:
	evt = &EnvVarTable{
		mux: new(sync.Mutex),
		et:  make(map[string]string),
	}
	// load and/or set default environment variables:
	evt.init()
	// set up hotreload watchdog:
	rwatch = &ReloadWatchdog{}
	rwatch.init(evt)
}

func main() {
	kubecontext, err := preflight()
	if err != nil {
		warn("Encountered issues during startup: " + err.Error())
	}
	rl, err = readline.NewEx(&readline.Config{
		AutoComplete:    completer,
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = rl.Close()
	}()
	setprompt(kubecontext)
	log.SetOutput(rl.Stderr())
	go rwatch.run()
	interpret(rl)
}

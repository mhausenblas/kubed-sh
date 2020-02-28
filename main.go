package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chzyer/readline"
	log "github.com/sirupsen/logrus"
)

const exitcmd = "exit"

var (
	version       string
	debugmode     bool
	prepull       bool
	customkubectl string
	prevdir       string
	rl            *readline.Instance
	completer     *readline.PrefixCompleter
	kplugins      map[string]string
)

func init() {
	if env := os.Getenv("KUBEDSH_DEBUG"); env != "" {
		debugmode = true
	}
	if env := os.Getenv("KUBEDSH_PREPULL"); env != "" {
		prepull = true
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
}

func main() {
	var script string
	// first, check if we've got a script filename
	// passed in via command line argument:
	if len(os.Args) == 2 {
		scriptfile := os.Args[1]
		b, err := ioutil.ReadFile(scriptfile)
		if err != nil {
			warn("Error executing script: " + err.Error())
		}
		script = string(b)
		interprets(script)
		return
	}
	// now let's see if we maybe have a script via stdin:
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			script += scanner.Text() + "\n"
		}
		if scanner.Err() != nil {
			warn("Error reading from stdin: " + scanner.Err().Error())
		}
		interprets(script)
		return
	}
	// well seems we're gonna be running interactive:
	err := preflight()
	if err != nil {
		warn("Encountered issues during startup: " + err.Error())
	}
	// set up auto-completion:
	autocompleter()
	defer func() {
		_ = rl.Close()
	}()
	// create and select global environment
	createenv(globalEnv, false)
	err = selectenv(globalEnv, false)
	if err != nil {
		warn("Encountered issues during startup: " + err.Error())
	}
	log.SetOutput(rl.Stderr())
	output("\nType 'help' to learn about available built-in commands.")
	// set up hotreload watchdog:
	rwatch = &ReloadWatchdog{}
	rwatch.init(currentenv().evt)
	go rwatch.run()
	// make jump pod available:
	go jpod()
	// necessary hack to make readline ignore a cascaded CTRL+C:
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			debug("caught an cascaded CTRL+C, ignoring it")
		}
	}()
	// kick off main interactive interpreter loop:
	interpreti(rl)
	// perform garbage collection on exit:
	doGC()
}

func jpod() {
	res, err := kubectl(false, "get", "po", "curljump")
	if strings.Contains(res, "NotFound") {
		return
	}
	_, err = kubectl(false, "run", "curljump", "--restart=Never",
		"--image=quay.io/mhausenblas/jump:0.2", "--", "sh", "-c", "sleep 10000")
	if err != nil {
		debug(err.Error())
	}
}

func autocompleter() {
	r, err := readline.NewEx(&readline.Config{
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
	})
	completer =
		readline.NewPrefixCompleter(
			readline.PcItem("cat"),
			readline.PcItem("cd"),
			readline.PcItem("curl"),
			readline.PcItem("cx"),
			readline.PcItem("echo"),
			readline.PcItem("env",
				readline.PcItem("list"),
				readline.PcItem("create"),
				readline.PcItem("select"),
				readline.PcItem("delete")),
			readline.PcItem(exitcmd),
			readline.PcItem("help"),
			readline.PcItem("img"),
			readline.PcItem("kill"),
			readline.PcItem("literally"),
			readline.PcItem("ls"),
			readline.PcItem("ns"),
			readline.PcItem("ps", readline.PcItem("all")),
			readline.PcItem("plugin", readline.PcItemDynamic(func(name string) []string {
				//auto-discover kubectl plugins and make them available natively:
				res, err := kubectl(false, "plugin", "list")
				if err != nil {
					fmt.Println(err.Error())
				}
				plugins := strings.Split(res, "\n")
				kplugins := make(map[string]string)
				for _, p := range plugins {
					cmd := filepath.Base(p)
					kplugins[strings.TrimPrefix(cmd, "kubectl-")] = p
				}
				plugincmds := make([]string, 0, len(kplugins))
				for k := range kplugins {
					plugincmds = append(plugincmds, k)
				}
				return plugincmds
			})),
			readline.PcItem("pwd"),
			readline.PcItem("sleep"),
			readline.PcItem("version"),
		)

	if err != nil {
		warn("Encountered issues during startup: " + err.Error())
	}
	r.Config.AutoComplete = completer
	rl = r
}

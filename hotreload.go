package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ReloadWatchdog watches for changes and updates binaries or scripts .
type ReloadWatchdog struct {
	active  bool
	et      *EnvVarTable
	watcher *fsnotify.Watcher
}

var (
	rwatch *ReloadWatchdog
)

func (rw *ReloadWatchdog) init(et *EnvVarTable) {
	rw.et = et
}

func (rw *ReloadWatchdog) isactive() bool {
	return rw.active
}

func (rw *ReloadWatchdog) checkstatus() {
	hr := rw.et.get("HOTRELOAD")
	hr = strings.ToLower(hr)
	rw.active = false
	if hr == "true" {
		rw.active = true
	}
}

func (rw *ReloadWatchdog) run() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		warn(err.Error())
	}
	rw.watcher = w
	defer func() {
		_ = rw.watcher.Close()
	}()
	done := make(chan bool)
	go rw.watch()
	err = rw.watcher.Add(".")
	if err != nil {
		warn(err.Error())
	}
	<-done
	debug("RUN DONE")
}

func (rw *ReloadWatchdog) watch() {
	for {
		rw.checkstatus()
		if rw.isactive() {
			select {
			case event := <-rw.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					debug("detected modify operation on " + event.String())
					f := strings.Split(event.Name, "!")[len(strings.Split(event.Name, "!"))-1]
					go rw.update(f)
				}
			case errw := <-rw.watcher.Errors:
				warn(errw.Error())
			}
		}
	}
}

func (rw *ReloadWatchdog) update(targetfile string) {
	// get rid of any silly stuff:
	targetfile = filepath.Clean(targetfile)
	debug("Restarting: " + targetfile)
	// find target pod and original file:
	po, err := kubectl(true, "get", "po",
		"--selector=script="+targetfile, "-o=custom-columns=:metadata.name", "--no-headers")
	if err != nil {
		debug(err.Error())
	}
	debug("updating pod " + po)
	res, err := kubectl(true, "get", "po", po,
		"-o=custom-columns=:metadata.annotations.original,:metadata.annotations.interpreter", "--no-headers")
	if err != nil {
		debug(err.Error())
	}
	original, interpreter := strings.Split(res, " ")[0], strings.Split(res, " ")[3]
	original = strings.TrimSpace(original)
	interpreter = strings.TrimSpace(interpreter)
	debug("original: " + original + " interpreter: " + interpreter)
	// copy changed file
	dest := fmt.Sprintf("%s:/tmp/", po)
	_, err = kubectl(false, "cp", original, dest)
	if err != nil {
		debug(err.Error())
	}
	// // kill in container
	_, _ = kubectl(false, "exec", po, "killall", interpreter)
	time.Sleep(5 * time.Second) // hack, right!?
	// start in container
	execremotescript := fmt.Sprintf("/tmp/%s", targetfile)
	_, err = kubectl(false, "exec", po, interpreter, execremotescript)
	if err != nil {
		debug(err.Error())
	}
	debug(targetfile + " updated in " + po)
}

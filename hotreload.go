package main

import (
	"fmt"
	"strings"

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
	fileq := make(chan string)
	go rw.queue(fileq)
	go rw.update(fileq)
	err = rw.watcher.Add(".")
	if err != nil {
		warn(err.Error())
	}
	<-done
}

func (rw *ReloadWatchdog) queue(fileq chan string) {
	for {
		rw.checkstatus()
		if rw.isactive() {
			select {
			case event := <-rw.watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					fileq <- event.Name
				}
			case errw := <-rw.watcher.Errors:
				warn(errw.Error())
			}
		}
	}
}

func (rw *ReloadWatchdog) update(fileq chan string) {
	for {
		targetfile := <-fileq
		info("Restarting: " + targetfile)
		// find target pod and original file
		po, err := kubectl(true, "get", "po",
			"--selector=script="+targetfile, "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			debug(err.Error())
			return
		}
		debug("updating pod " + po)
		res, err := kubectl(true, "get", "po", po,
			"-o=custom-columns=:metadata.annotations.original,:metadata.annotations.interpreter", "--no-headers")
		if err != nil {
			debug(err.Error())
			return
		}
		debug("result of annotations query: " + res)
		original, interpreter := strings.Split(res, " ")[0], strings.Split(res, " ")[1]
		original = strings.TrimSpace(original)
		interpreter = strings.TrimSpace(interpreter)
		// copy changed file
		dest := fmt.Sprintf("%s:/tmp/", po)
		_, err = kubectl(true, "cp", original, dest)
		if err != nil {
			debug(err.Error())
			return
		}
		// kill in container
		_, _ = kubectl(true, "exec", po, "--", "killall", "-9", interpreter)
		if err != nil {
			debug(err.Error())
			return
		}
		// start in container
		execremotescript := fmt.Sprintf("/tmp/%s", targetfile)
		res, err = kubectl(true, "exec", po, "--", interpreter, execremotescript)
		if err != nil {
			debug(err.Error())
			return
		}
		debug("result of launching new file: " + res)
	}
}

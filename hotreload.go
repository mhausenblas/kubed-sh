package main

import (
	"strings"
	"time"
)

// ReloadWatchdog watches for changes and updates binaries or scripts .
type ReloadWatchdog struct {
	active bool
	et     *EnvVarTable
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
	if hr == "true" {
		rw.active = true
	}
}

func (rw *ReloadWatchdog) run() {
	for {
		rw.checkstatus()
		if rw.isactive() {
			info("Checking for changes")
		}
		time.Sleep(1 * time.Second)
	}
}

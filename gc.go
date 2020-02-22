package main

import (
	"strings"
	"time"
)

var (
	// how often we check for orphans:
	gcPause = 30 * time.Second
)

func gcDProcs() {
	for {
		orphandeploys, err := kubectl(false, "get", "deploy",
			"--selector=dproctype="+string(DProcTerminating), "-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			debug(err.Error())
		}
		debug(orphandeploys)
		if orphandeploys != "" {
			for _, d := range strings.Split(orphandeploys, "\n") {
				_, err := kubectl(false, "delete", "deploy", d)
				if err != nil {
					warn("GC: couldn't reap orphaned deployment " + d)
				}
			}
		}
		time.Sleep(gcPause)
	}
}

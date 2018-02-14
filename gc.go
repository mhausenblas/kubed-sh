package main

import (
	"strings"
	"time"
)

var (
	// how often we check for orphans:
	gcPause = 10 * time.Second
	// how long a pod of a terminating dproc can run
	// before we consider it an orphan (in seconds):
	maxOrphanRuntimeSec = 600.0
)

func gcDProcs() {
	for {
		poNstart, err := kubectl(false, "get", "po",
			"--selector=dproctype="+string(DProcTerminating), "-o=custom-columns=:metadata.name,:status.startTime", "--field-selector=status.phase=Running",
			"--no-headers")
		if err != nil {
			debug(err.Error())
		}
		debug(poNstart)
		if poNstart != "" {
			for _, pns := range strings.Split(poNstart, "\n") {
				poname, start := strings.Split(pns, "   ")[0], strings.Split(pns, "   ")[1]
				debug("GC: looking at candidate pod " + poname + " with start timestamp " + start + "\n")
				layout := "2006-01-02T15:04:05Z"
				st, err := time.Parse(layout, start)
				if err != nil {
					debug("GC: couldn't parse start time of pod " + poname)
				}
				now := time.Now()
				diff := now.Sub(st)
				if diff.Seconds() > maxOrphanRuntimeSec {
					debug("GC: found orphaned pod " + poname + " with a start time of " + st.String())
					_, err = kubectl(false, "delete", "pod", poname)
					if err != nil {
						warn("GC: couldn't remove orphaned pod " + poname)
					}
				}
			}
		}
		time.Sleep(gcPause)
	}
}

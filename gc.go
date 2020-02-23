package main

import (
	"strings"
)

func gcDProcs() {
	// reap jump pod:
	_, err := kubectl(false, "delete", "po", "curljump")
	if err != nil {
		warn("GC: couldn't reap jump pod")
	}
	// reap orphaned pods:
	orphandpods, err := kubectl(false, "get", "po",
		"--selector=gen=kubed-sh", "-o=custom-columns=:metadata.name", "--no-headers")
	if err != nil {
		debug(err.Error())
	}
	debug(orphandpods)
	if orphandpods != "" {
		for _, p := range strings.Split(orphandpods, "\n") {
			_, err := kubectl(false, "delete", "po", p)
			if err != nil {
				warn("GC: couldn't reap orphaned pod " + p)
			}
		}
	}
	// reap orphaned deployments:
	orphandeploys, err := kubectl(false, "get", "deploy",
		"--selector=gen=kubed-sh", "-o=custom-columns=:metadata.name", "--no-headers")
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
	// reap orphaned services:
	orphandesvcs, err := kubectl(false, "get", "services",
		"--selector=gen=kubed-sh", "-o=custom-columns=:metadata.name", "--no-headers")
	if err != nil {
		debug(err.Error())
	}
	debug(orphandesvcs)
	if orphandesvcs != "" {
		for _, s := range strings.Split(orphandesvcs, "\n") {
			_, err := kubectl(false, "delete", "service", s)
			if err != nil {
				warn("GC: couldn't reap orphaned service " + s)
			}
		}
	}
}

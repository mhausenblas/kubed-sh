package main

import (
	"fmt"
	"os"
	"strings"
)

func doGC() {
	gcoptions := []string{""}
	val, ok := os.LookupEnv("KUBEDSH_GC")
	if ok {
		switch {
		case strings.Contains(val, ","):
			gcoptions = strings.Split(val, ",")
		default:
			gcoptions[0] = val
		}
		info(fmt.Sprintf("Performing garbage collection for %v, this may take a few seconds ...", gcoptions))
		for _, gco := range gcoptions {
			switch gco {
			case "JUMP_POD": // reap the jump pod:
				_, err := kubectl(false, "delete", "po", "curljump")
				if err != nil {
					warn("GC: couldn't reap jump pod")
				}
			case "ALL_PODS": // reap all orphaned pods:
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
			case "ALL_DEPLOYS": // reap all orphaned deployments:
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
			case "ALL_SVC": // reap all orphaned services:
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
		}
		return
	}
	info("Skipping garbage collection, no strategy selected. Set via the KUBEDSH_GC environment variable.")

}

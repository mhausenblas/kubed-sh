package main

import (
	"fmt"
	"os"
	"sync"
)

// EnvVarTable is an environment variable table.
type EnvVarTable struct {
	mux *sync.Mutex
	et  map[string]string
}

// Environment represents an environment.
type Environment struct {
	name string
	evt  *EnvVarTable
}

var (
	globalEnv    = "KUBED-SH_GLOBAL_ENVIRONMENT"
	currentEnv   *Environment
	environments = make(map[string]Environment)
)

func currentenv() *Environment {
	return currentEnv
}

func createenv(name string, showhint bool) {
	// set up the environment variables table for our new environment:
	evt := &EnvVarTable{
		mux: new(sync.Mutex),
		et:  make(map[string]string),
	}
	// load and/or set default environment variables for the environment variables table:
	evt.init()
	env := Environment{
		name: name,
		evt:  evt,
	}
	environments[name] = env
	if showhint {
		info("Created environment [" + name + "]")
		info("To activate it, use env select " + name)
	}
}

func selectenv(name string, showhint bool) error {
	env, ok := environments[name]
	if !ok {
		return fmt.Errorf("provided environment doesn't seem to exist")
	}
	currentEnv = &env
	setprompt()
	if showhint && currentenv().name != globalEnv {
		info("Selected environment [" + name + "]")
	}
	return nil
}

func deleteenv(name string, showhint bool) error {
	_, ok := environments[name]
	if !ok {
		return fmt.Errorf("provided environment doesn't seem to exist")
	}
	// change the environment of all dprocs in the environment:
	for dpid, dproc := range dpt.lt {
		if dproc.Env == name {
			dproc.Env = globalEnv
			dpt.lt[dpid] = dproc
		}
	}
	// re-label the resources to global env:
	_, err := kubectl(false, "label", "deploy,rs,svc,po", "--selector=env="+name, "--overwrite", "env="+globalEnv)
	if err != nil {
		warn("Can't move processes from " + name + " to global environment")
	}
	if showhint {
		info("Deleted environment [" + name + "], now all processes are in the global environment")
	}
	// set current env to global env and get rid of env
	_ = selectenv(globalEnv, true)
	delete(environments, name)
	return nil
}

// set sets an environment variable
func (et *EnvVarTable) set(envar, value string) {
	et.mux.Lock()
	et.et[envar] = value
	et.mux.Unlock()
}

// get returns the value of an environment variable
func (et *EnvVarTable) get(envar string) string {
	et.mux.Lock()
	val, ok := et.et[envar]
	et.mux.Unlock()
	if !ok {
		return ""
	}
	return val
}

// unset removes an environment variable
func (et *EnvVarTable) unset(envar string) {
	et.mux.Lock()
	delete(et.et, envar)
	et.mux.Unlock()
}

// init sets default env vars and loads some
// such as $PATH, $HOME, etc. from parent shell.
func (et *EnvVarTable) init() {
	// set defaults:
	et.set("SERVICE_PORT", "80")
	et.set("SERVICE_NAME", "")
	et.set("BINARY_IMAGE", "alpine:3.7")
	et.set("NODE_IMAGE", "node:9.4-alpine")
	et.set("PYTHON_IMAGE", "python:3.6-alpine3.7")
	et.set("RUBY_IMAGE", "ruby:2.5-alpine3.7")
	et.set("HOTRELOAD", "false")
	// load from parent shell, if present:
	val, ok := os.LookupEnv("KUBECTL_BINARY")
	if ok {
		et.set("KUBECTL_BINARY", val)
	}
	val, ok = os.LookupEnv("PATH")
	if ok {
		et.set("PATH", val)
	}
	val, ok = os.LookupEnv("HOME")
	if ok {
		et.set("HOME", val)
	}
}

func setprompt() {
	context, err := kubectl(false, "config", "current-context")
	if err != nil {
		warn("Can't determine current context")
	}
	namespace, err := kubectl(false, "run", "ns", "--rm", "-i", "-t", "--restart=Never", "--image=alpine:3.7", "--", "cat", "/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		warn("Can't determine namespace")
	}
	env := currentenv().name
	switch env {
	case globalEnv:
		rl.SetPrompt(fmt.Sprintf("[\033[32m%s\033[0m::\033[36m%s\033[0m]$ ", context, namespace))
	default:
		rl.SetPrompt(fmt.Sprintf("[\033[95m%s\033[0m@\033[32m%s\033[0m::\033[36m%s\033[0m]$ ", env, context, namespace))
	}
}

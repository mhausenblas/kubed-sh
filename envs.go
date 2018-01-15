package main

import "sync"

// EnvVarTable is the global environment variable table.
type EnvVarTable struct {
	mux *sync.Mutex
	et  map[string]string
}

var (
	evt *EnvVarTable
)

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

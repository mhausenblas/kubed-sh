package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"sync"
	"text/tabwriter"
)

// DProcType represents the type of distributed process.
type DProcType string

// DProc represents a distributed process.
// A distributed process runs in the Kubernetes
// cluster in a certain context. It has an owner
// (either a human or a deployment) and the source
// (the location of either the binary or the script).
// All DProcs are labelled with gen=kubed-sh and,
// depending on the type of execution with either
// script=xxx or bin=xxx.
type DProc struct {
	ID          string
	Type        DProcType
	KubeContext string
	Src         string
}

// DProcTable is the distributed process lookup table,
// mappping
type DProcTable struct {
	mux *sync.Mutex
	lt  map[string]DProc
}

var (
	// DProcTerminating stands for a terminating distributed process
	// (one-shot, batch) which will be launched via a pod
	DProcTerminating DProcType = "terminating"
	// DProcLongRunning stands for a long-running distributed process
	// which will be launched via a deployment + service
	DProcLongRunning DProcType = "longrunning"
	// The global distributed process lookup table.
	dpt *DProcTable
)

// BuildDPT adds a distributed process to the global table.
func (dt *DProcTable) BuildDPT() error {
	kubecontext, err := kubectl("config", "current-context")
	if err != nil {
		return err
	}
	res, err := kubectl("get", "deployments", "--selector=gen=kubed-sh",
		"-o=custom-columns=:metadata.name,:metadata.labels", "--no-headers")
	if err != nil {
		return fmt.Errorf("Failed to gather distributed processes due to:\n%s", err)
	}
	if res == "" {
		return nil
	}
	debug(res)
	for _, r := range strings.Split(res, "\n") {
		// now r is something like 'kubed-sh-1516013421817997000   map[gen:kubed-sh script:test.js]'
		id := strings.Split(r, "   ")[0]
		labels := strings.Split(r, "   ")[1]
		// now labels is something like map[gen:kubed-sh script:test.js]
		labels = strings.TrimSuffix(labels[4:], "]")
		// now labels is something like gen:kubed-sh script:test.js
		src := ""
		for _, s := range strings.Split(labels, " ") {
			if strings.HasPrefix(s, "script") || strings.HasPrefix(s, "bin") {
				src = s
				break
			}
		}
		debug("id: " + id + " source: " + src)
		dt.addDProc(newDProc(id, DProcLongRunning, kubecontext, src))
	}
	return nil
}

// DumpDPT lists all distributed processes.
func (dt *DProcTable) DumpDPT(kubecontext string) string {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DPID\tCONTEXT\tSOURCE\t")
	for _, dproc := range dt.lt {
		switch kubecontext {
		case "":
			fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t", dproc.ID, dproc.KubeContext, dproc.Src))
		case dproc.KubeContext:
			fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t", dproc.ID, dproc.KubeContext, dproc.Src))
		}
	}
	_ = tw.Flush()
	_ = w.Flush()
	return b.String()
}

// newDProc creates a distributed process entry.
func newDProc(dpid string, dptype DProcType, context, source string) DProc {
	return DProc{
		ID:          dpid,
		Type:        dptype,
		KubeContext: context,
		Src:         source,
	}
}

func (dproc DProc) String() string {
	return fmt.Sprintf("%v %v %v %v", dproc.ID, dproc.Type, dproc.KubeContext, dproc.Src)
}

// addDProc adds a distributed process to the global table.
func (dt *DProcTable) addDProc(dproc DProc) {
	// the lookup key is composed of the distributed process and the context
	k := fmt.Sprintf("%s@%s", dproc.ID, dproc.KubeContext)
	dt.mux.Lock()
	dt.lt[k] = dproc
	dt.mux.Unlock()
}

// removeDProc removes a distributed process from the global table.
func (dt *DProcTable) removeDProc(dproc DProc) {
	// the lookup key is composed of the distributed process and the context
	k := fmt.Sprintf("%s@%s", dproc.ID, dproc.KubeContext)
	dt.mux.Lock()
	delete(dt.lt, k)
	dt.mux.Unlock()
}

// getDProc looks up a distributed process in the global table
// based on the distributed process ID and the context.
func (dt *DProcTable) getDProc(ID, context string) (DProc, error) {
	// the lookup key is composed of the distributed process and the context
	k := fmt.Sprintf("%s@%s", ID, context)
	dt.mux.Lock()
	d, ok := dt.lt[k]
	dt.mux.Unlock()
	if !ok {
		return d, fmt.Errorf("Distributed process %s not found", k)
	}
	return d, nil
}

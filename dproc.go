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
	// TERMINATING stands for a terminating distributed process
	// (one-shot, batch) which will be launched via a pod
	TERMINATING DProcType = "terminating"
	// LONGRUNNING stands for a long-running distributed process
	// which will be launched via a deployment + service
	LONGRUNNING DProcType = "longrunning"
	// The global distributed process lookup table.
	dpt *DProcTable
)

// BuildDPT adds a distributed process to the global table.
func (dt *DProcTable) BuildDPT() error {
	res, err := kubectl("get", "pods,deployments", "--selector=gen=kubed-sh",
		"-o=custom-columns=:metadata.name", "--no-headers")
	if err != nil {
		return fmt.Errorf("Failed to gather distributed processes due to:\n%s", err)
	}
	kubecontext, err := kubectl("config", "current-context")
	if err != nil {
		return err
	}
	info(res)
	for _, id := range strings.Split(res, "\n") {
		dt.addDProc(newDProc(id, LONGRUNNING, kubecontext, "source.example"))
	}
	return nil
}

// DumpDPT lists all distributed processes.
func (dt *DProcTable) DumpDPT() string {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DPID\tCONTEXT\tSOURCE\t")
	for _, dproc := range dt.lt {
		fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t", dproc.ID, dproc.KubeContext, dproc.Src))
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

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
	ServiceName string
	Env         string
}

// DProcTable is the distributed process lookup table,
// mappping distributed process IDs (dpids) to distributed
// processes (dproc)
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
	kubecontext, err := kubectl(true, "config", "current-context")
	if err != nil {
		return err
	}
	res, err := kubectl(true, "get", "deployments", "--selector=gen=kubed-sh",
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
		// now labels is something like map[gen:kubed-sh script:test.js env:abc]
		labels = strings.TrimSuffix(labels[4:], "]")
		// now labels is something like gen:kubed-sh script:test.js env:abc
		var src, env string
		// now grab source and env labels:
		for _, s := range strings.Split(labels, " ") {
			if strings.HasPrefix(s, "script") || strings.HasPrefix(s, "bin") {
				src = s
			}
			if strings.HasPrefix(s, "env") {
				env = strings.Split(s, ":")[1]
				createenv(env, false)
			}
		}
		debug("env:" + env + " id: " + id + " source: " + src)
		svcname, err := kubectl(true, "get", "services", "--selector=gen=kubed-sh,"+strings.Replace(src, ":", "=", -1),
			"-o=custom-columns=:metadata.name", "--no-headers")
		if err != nil {
			return fmt.Errorf("Failed to gather distributed processes due to:\n%s", err)
		}
		dt.addDProc(newDProc(id, DProcLongRunning, kubecontext, src, svcname, env))
	}
	return nil
}

// DumpDPT lists all distributed processes.
func (dt *DProcTable) DumpDPT(kubecontext string) string {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	switch kubecontext {
	case "":
		fmt.Fprintln(tw, "DPID\tENV\tCONTEXT\tSOURCE\tURL")
		for _, dproc := range dt.lt {
			switch dproc.Env {
			case globalEnv:
				fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t", dproc.ID, "global", dproc.KubeContext, strings.Split(dproc.Src, ":")[1], dproc.ServiceName))
			default:
				fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t", dproc.ID, dproc.Env, dproc.KubeContext, strings.Split(dproc.Src, ":")[1], dproc.ServiceName))
			}

		}
	default:
		fmt.Fprintln(tw, "DPID\tSOURCE\tURL")
		for _, dproc := range dt.lt {
			if dproc.KubeContext == kubecontext && dproc.Env == currentenv().name {
				fmt.Fprintln(tw, fmt.Sprintf("%s\t%s\t%s\t", dproc.ID, strings.Split(dproc.Src, ":")[1], dproc.ServiceName))
			}
		}
	}
	_ = tw.Flush()
	_ = w.Flush()
	return b.String()
}

// newDProc creates a distributed process entry.
func newDProc(dpid string, dptype DProcType, context, source, svcname, env string) DProc {
	return DProc{
		ID:          dpid,
		Type:        dptype,
		KubeContext: context,
		Src:         source,
		ServiceName: svcname,
		Env:         env,
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

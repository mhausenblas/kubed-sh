package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func hsleep(line string) {
	if !strings.ContainsAny(line, " ") {
		warn("Need a time interval in seconds. For example, `sleep 10s` means sleep for 10s and `sleep 200ms` means sleep for 200ms")
		return
	}
	arg := strings.Split(line, " ")[1]
	d, err := time.ParseDuration(arg)
	if err != nil {
		warn("Can't parse time interval:" + err.Error())
		return
	}
	time.Sleep(d)
}

func hcd(line string) {
	var targetdir string
	switch {
	case !strings.ContainsAny(line, " "):
		usr, err := user.Current()
		if err != nil {
			warn(err.Error())
		}
		targetdir = usr.HomeDir
	case strings.ContainsAny(line, " ") && strings.Split(line, " ")[1] == "-":
		targetdir = prevdir
	default:
		targetdir = strings.Split(line, " ")[1]
	}
	prevdir, _ = os.Getwd()
	err := os.Chdir(targetdir)
	if err != nil {
		warn(err.Error())
	}
}

func hcurl(line string) {
	if !strings.ContainsAny(line, " ") {
		info("Need a target URL, for example `curl someservice` in the cluster or `curl http://example.com`")
		return
	}
	url := strings.Split(line, " ")[1]
	res, err := kubectl(false, "exec", "-it", "curljump", "curl", url)
	if err != nil {
		warn(fmt.Sprintf("Can't curl %s: %s", url, err.Error()))
		return
	}
	output(res)
}

func hlocalexec(line string) {
	cmd := line
	args := []string{}
	if strings.ContainsAny(line, " ") {
		cmd = strings.Split(line, " ")[0]
		args = strings.Split(line, " ")[1:]
	}
	res, err := shellout(true, cmd, args...)
	if err != nil {
		fmt.Printf("Can't execute %s locally: %s", cmd, err)
	}
	output(res)
}

func henv(line string) {
	// user asked to list variables in the currently selected environment:
	if !strings.ContainsAny(line, " ") {
		tmp := []string{}
		for k, v := range currentenv().evt.et {
			tmp = append(tmp, fmt.Sprintf("%s=%s", k, v))
		}
		sort.Strings(tmp)
		for _, e := range tmp {
			output(e)
		}
		return
	}
	// maybe it's a list command?
	if line == "env list" {
		tmp := []string{}
		for k := range environments {
			if k != globalEnv {
				tmp = append(tmp, fmt.Sprintf("%s", k))
			}
		}
		sort.Strings(tmp)
		for _, e := range tmp {
			output(e)
		}
		return
	}
	// so check for a CRUD command?
	if (len(strings.Split(line, " "))) != 3 {
		warn("Unknown command. Must follow 'env list' or 'env create|select|delete ENV_NAME' pattern.")
		return
	}
	cmd := strings.Split(line, " ")[1]
	targetenv := strings.Split(line, " ")[2]
	switch cmd {
	case "create":
		debug("creating new environment '" + targetenv + "'")
		createenv(targetenv, true)
	case "select":
		debug("switching to environment '" + targetenv + "'")
		err := selectenv(targetenv, true)
		if err != nil {
			warn(err.Error())
		}
	case "delete":
		debug("deleting environment '" + targetenv + "'")
		err := deleteenv(targetenv, true)
		if err != nil {
			warn(err.Error())
		}
	default:
		warn("Unknown command. Must follow 'env list' or 'env create|select|delete ENV_NAME' pattern.")
	}
}

func hkill(line string) {
	if !strings.ContainsAny(line, " ") {
		info("Need a target distributed process to kill")
		return
	}
	// pre-flight check if dproc exists:
	ID := strings.Split(line, " ")[1]
	_, err := kubectl(true, "get", "deployment", ID)
	if err != nil {
		warn(fmt.Sprintf("A distributed process with the ID '%s' does not exist in current context. Try the ps command …", ID))
		return
	}
	_, err = kubectl(true, "scale", "--replicas=0", "deployment", ID)
	if err != nil {
		killfail(line, err.Error())
		return
	}
	_, err = kubectl(true, "delete", "deployment", ID)
	if err != nil {
		killfail(line, err.Error())
		return
	}

	// gather info to remove from global DPT:
	kubecontext, err := kubectl(true, "config", "current-context")
	if err != nil {
		killfail(line, err.Error())
		return
	}
	dproc, err := dpt.getDProc(ID, kubecontext)
	if err != nil {
		killfail(line, err.Error())
		return
	}
	svcname := dproc.ServiceName
	debug(svcname)
	_, err = kubectl(true, "delete", "service", svcname)
	if err != nil {
		killfail(line, err.Error())
		return
	}
	if err != nil {
		info(err.Error())
	}
	// finally, remove drpoc from global DPT:
	dpt.removeDProc(dproc)
}

func hps(line string) {
	args := ""
	if strings.ContainsAny(line, " ") {
		args = strings.Split(line, " ")[1]
	}
	var kubecontext string
	switch args {
	case "all":
		kubecontext = ""
	default:
		k, err := kubectl(true, "config", "current-context")
		if err != nil {
			warn("Can't determine current context")
			return
		}
		kubecontext = k
	}
	debug("in context " + kubecontext)
	res := dpt.DumpDPT(kubecontext)
	output(res)
}

func hcontexts(line string) {
	if !strings.ContainsAny(line, " ") {
		res, err := kubectl(true, "config", "get-contexts")
		if err != nil {
			warn(fmt.Sprintf("Can't list contexts: %s", err))
			return
		}
		output(res)
		return
	}
	targetcontext := strings.Split(line, " ")[1]
	res, err := kubectl(true, "config", "use-context", targetcontext)
	if err != nil {
		warn(fmt.Sprintf("Can't switch context: %s", err))
		return
	}
	output(res)
	if rl != nil {
		setprompt()
	}
}

func hns(line string) {
	if !strings.ContainsAny(line, " ") {
		res, err := kubectl(true, "get", "ns")
		if err != nil {
			warn(fmt.Sprintf("Can't list namespaces: %s", err))
			return
		}
		output(res)
		return
	}
	currentcx, err := kubectl(false, "config", "current-context")
	if err != nil {
		warn("Can't determine current context")
		return
	}
	targetns := strings.Split(line, " ")[1]
	res, err := kubectl(true, "config", "set-context", currentcx,
		"--namespace="+targetns)
	if err != nil {
		warn(fmt.Sprintf("Can't switch namespace: %s", err))
		return
	}
	output(res)
	if rl != nil {
		setprompt()
	}
}

func himg() {
	res, err := kubectl(true, "get", "po", "--all-namespaces", "-o", "jsonpath=\"{..image}\"")
	if err != nil {
		warn(fmt.Sprintf("Can't list container images: %s", err))
		return
	}
	seenimgs := []string{}
	for _, img := range strings.Split(strings.Trim(res, "\""), " ") {
		if !contains(seenimgs, img) {
			seenimgs = append(seenimgs, img)
		}
	}
	sort.Strings(seenimgs)
	for _, uimg := range seenimgs {
		output(uimg)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func launchfail(line, reason string) {
	fmt.Printf("\nCan't launch %s in the cluster:\n%s\n\n", strconv.Quote(line), reason)
}

func killfail(line, reason string) {
	fmt.Printf("\nCan't kill %s:\n%s\n\n", strconv.Quote(line), reason)
}

func extractsrc(line string) string {
	line = strings.TrimSuffix(line, "&")
	line = strings.TrimSpace(line)
	debug("input line: " + line)
	// a binary is standalone:
	if !strings.ContainsAny(line, " ") {
		_, binfile := filepath.Split(line)
		debug("binary extracted: " + binfile)
		return binfile
	}
	// … otherwise it's a script:
	script := strings.Split(line, " ")[1]
	_, scriptfile := filepath.Split(script)
	debug("script extracted: " + scriptfile)
	return scriptfile
}

func hlaunch(line string) {
	// If a line doesn't start with one of the
	// known environments, assume user wants to
	// launch a binary:
	var dpid, svcname string
	src := extractsrc(line)
	src = "script:" + src
	switch {
	case strings.HasPrefix(line, "python "):
		d, s, err := launchpy(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
		svcname = s
	case strings.HasPrefix(line, "node "):
		d, s, err := launchjs(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
		svcname = s
	case strings.HasPrefix(line, "ruby "):
		d, s, err := launchrb(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
		svcname = s
	default: // binary
		d, s, err := launchbin(line)
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpid = d
		svcname = s
		src = "bin:" + extractsrc(line)
	}
	// update DPT
	if strings.HasSuffix(line, "&") {
		kubecontext, err := kubectl(true, "config", "current-context")
		if err != nil {
			launchfail(line, err.Error())
			return
		}
		dpt.addDProc(newDProc(dpid, DProcLongRunning, kubecontext, src, svcname, currentenv().evt.get("SERVICE_PORT"), currentenv().name))
	}
}

func hliterally(line string) {
	if !strings.ContainsAny(line, " ") {
		info("Can't execute kubectl command, not enough arguments")
		return
	}
	l := strings.Split(line, " ")
	res, _ := kubectl(true, l[1], l[2:]...)
	output(res)
}

func hecho(line string) {
	if !strings.ContainsAny(line, " ") {
		info("No value to echo given")
		return
	}
	echo := strings.Split(line, " ")[1]
	if strings.HasPrefix(echo, "$") {
		v := currentenv().evt.get(echo[1:])
		if v != "" {
			fmt.Println(v)
			return
		}
	}
	fmt.Println(echo)
}

func hplugins(line string) {
	cmd := line
	args := []string{}
	if strings.ContainsAny(line, " ") {
		cmd = strings.Split(line, " ")[1]
		args = strings.Split(line, " ")[2:]
	}
	res, err := kubectl(true, cmd, args...)
	if err != nil {
		fmt.Printf("Can't execute %s locally: %s", cmd, err)
	}
	output(res)
}

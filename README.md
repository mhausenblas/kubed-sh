# kubed-sh

Hello and welcome to `kubed-sh`, the Kubernetes distributed shell for the casual cluster user.
If you have access to a [Kubernetes](https://kubernetes.io/) cluster, you can [install it](#install-it) now
and then learn how to [use it](#use-it). In a nutshell `kubed-sh` ([pronunciation](#faq)) lets you execute
a program in a Kubernetes cluster without having to create a container image or learn new commands. For example:

![Launching a simple binary in a Kubernetes cluster](img/launch-bin.png)

Above you see the Linux ELF binary `tc` that you get when doing a `GOOS=linux go build` in the test case directory [tc/](tc/),
executing in the Kubernetes cluster, producing the output `I'm a simple program that just prints this message and exits`.

In addition to launching (Linux ELF) binaries, the following interpreted environments are currently supported:

- When you enter `node script.js`, a Node.js (default version: 9.4) environment is provided and `script.js` is executed.
- When you enter `python script.py`, a Python (default version: 3.6) environment is provided and the `script.py` is executed.
- When you enter `ruby script.rb`, a Ruby (default version: 2.5) environment is provided and the `script.rb` is executed.

Note that `kubed-sh` is a proper shell environment, that is, you can expect features such as auto-complete, history operations,
or `CTRL+L` clearing the screen to work as per usual. Also, you can read here [why](why.md) I wrote `kubed-sh`.

## Install it

Prerequisites:

1. You need [Go](https://golang.org/dl/) in order to build it as there are no binaries, currently. I'm using `go1.9.2 darwin/amd64` on my machine.
1. `kubectl` must be installed, I tested it with client version 1.9.1 so far.
1. Access to a Kubernetes cluster must be configured:
  - if you do `ls ~/.kube/config > /dev/null && echo $?` and you see a `0` as a result, you're good, and further
  - if you do `kubectl config get-contexts | wc -l` and see a number greater than `0`, then that's super dope.

Now to install `kubed-sh` simply do the following (anywhere):

```
$ go get github.com/mhausenblas/kubed-sh
```

Note that if your `$GOPATH/bin` is in your `$PATH` then now you can use `kubed-sh` from everywhere. If not, you can:

- Do a `cd $GOPATH/src/github.com/mhausenblas/kubed-sh` followed by a `go build` and use it from this directory.
- Run it like so: `$GOPATH/bin/kubed-sh`

## Use it

Once you've `kubed-sh` installed, launch it and you should find yourself in an interactive shell, that is, a [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) like so:

```
$ kubed-sh
Note: It seems you're running kubed-sh in a non-Linux environment (detected: darwin),
so make sure the binaries you launch are Linux binaries in ELF format.

[minikube]$
```

Supported commands (see also `help`):

- `contexts` (local) … list available Kubernetes contexts (cluster, namespace, user tuples)
- `echo` (local) … print a value or environment variable
- `env`(local) … list all environment variables currently defined
- `exit` (or: `quit`, local) … leave shell
- `help` (local) … list built-in commands
- `kill` (distributed) … stop a distributed process
- `literally` (or prefix with `` ` ``, local) … drop down to raw mode, literally execute as a kubectl command
- `ps` (distributed) … list all distributed (long-running) processes in current context
- `pwd` (local) … print current working directory
- `use` (local) … select a certain context to work with

## FAQ

**Q**: Why another Kubernetes shell? There's already [cloudnativelabs/kube-shell](https://github.com/cloudnativelabs/kube-shell) and [c-bata/kube-prompt](https://github.com/c-bata/kube-prompt). <br>
**A**: True, there is previous art, though these shells more or less aim at making `kubectl` interactive, exposing the commands such as `get` or `apply` to the user.
In a sense `kubed-sh` is more like [technosophos/kubeshell](https://github.com/technosophos/kubeshell), trying to provide an environment a typical *nix user is comfortable with.
For example, rather than providing a `create` or `apply` command to run a program, the user would simply enter the name of the executable, as she would do, for example, in the bash shell. See also the [motivation](why.md).

**Q**: How is `kubed-sh` pronounced? <br>
**A**: Well, I pronounce it /ku:bˈdæʃ/ as in 'kube dash' ;)

**Q**: How does this actually work? <br>
**A**: Good question. Essentially a glorified `kubectl` wrapper on steriods. See also the [design](design.md).

Once you've `kubed-sh` installed, launch it and you should find yourself in an interactive shell:

```sh
$ kubed-sh
Note: It seems you're running kubed-sh in a non-Linux environment (detected: darwin),
so make sure the binaries you launch are Linux binaries in ELF format.

Detected Kubernetes client in version v1.17.2 and server in version v1.14.9-eks-c0eccc

[user@somecluster.eu-west-1.eksctl.io::demo]$
```

Above, you notice that on start-up `kubed-sh` will tell you which client and 
server version of Kubernetes it has detected and at any point in time you are 
able to tell in which context you're working (here: `user@somecluster.eu-west-1.eksctl.io`) 
as well in which namespace (here: `demo`) you're operating. Learn [more …](../design/#environments)


## Built-ins

In general `kubed-sh` tries to behave like a normal/local shell such as Bash.
For example, if an input (or line in a script) starts with an `#` it's considered
a comment and hence ignored. There are however some cases where `kubed-sh` differs
significantly from a local shell: the execution of what we call distributed processes.
If you launch such a distributed process, either via entering the name of a Linux binary
or via one of the supported interpreted languages (Node.js, Python, Ruby), then
under the hood and transparently to you `kubed-sh` creates Kubernetes resources 
in the target cluster (namespace). Keep this in mind when you see the `scope` of 
a command; if it says `cluster` that's a hint that under the hood some `kubectl`
voodoo is going on.

### Commands

The built-in commands of `kubed-sh` (see also `help`) are as follows:

| command    | scope | description                                             |
| ----------:| ----- | ------------------------------------------------------- |
| `cat`      | local | output the content of file to terminal                  |
| `cd`       | local | change working directory                                |
| `curl`     | cluster | execute a curl operation in the cluster               |
| `cx`       | local | list/select Kubernetes contexts                         |
| `debug`    | local | toggle debug mode (show kubectl commands, etc.)         | 
| `echo`     | local | print a value or environment variable                   |
| `env`      | local | see below                                               |
| `exit`     | local | exit `kubed-sh`                                         |
| `help`     | local | list built-in commands; use help `command` for details  |
| `img`      | local | list container images                                   |
| `kill`     | cluster | stop a distributed process                            |
| `literally`| local | literally execute as a `kubectl` command                |
| `ls`       | local | list content of directory                               |
| `ns`       | local | list/select Kubernetes namespaces                       |
| `plugin`   | local | list/execute kubectl plugin                             |
| `ps`       | cluster |  list all distributed processes in current environment|
| `pwd`      | local | print current working directory                         |
| `sleep`    | local | sleep for specified time interval (a NOP)               |
| `version ` | local | print `kubed-sh` version                                |

!!! tip
      Rather than the lengthy `literally` simply prefix a line with \`
      to achieve the same. For example, if you enter ``get po` you list the pods
      in the current namespace.

The `env` command (which, on its own lists `kubed-sh` [environment variables](#environment-variables))
has four sub-commands:

- `env list` … list all defined environments in current context
- `env create $ENVNAME` … create a new environment called `$ENVNAME`
- `env select $ENVNAME` … make environment called `$ENVNAME` the active one
- `env delete $ENVNAME` … delete environment called `$ENVNAME`

If no environment is selected, you are operating in the global environment.

!!! note
      When you execute the `env delete ENVNAME` command, the environment is 
      reaped and all the distributed processes go back into the global environment.

### Environment variables

In `kubed-sh` you can define and use your own [environments](../design/#environments).
Within an environment, there are a number of pre-defined environment variables,
which influence the creation of the distributed processes:

- `BINARY_IMAGE` (default: `alpine:3.7`) … for executing binaries
- `NODE_IMAGE` (default: `node:12-alpine`) … for executing Node.js scripts
- `PYTHON_IMAGE` (default: `python:3.6-alpine3.7`) … for executing Python scripts
- `RUBY_IMAGE` (default: `ruby:2.5-alpine3.7`) … for executing Ruby scripts
- `SERVICE_PORT` (default: `80`) … expose long-running process on this port (in-cluster)
- `SERVICE_NAME` (default: `""`) … overwrite URL of long-running process (in-cluster)
- `HOTRELOAD` (default: `false`) … enable a watch on local files to trigger automatic updates on modification (EXPERIMENTAL)

!!! tip
      You can overwrite any of the above environment variables to change the 
      runtime behavior of the distributed processes you create. All changes are 
      valid for the runtime of `kubed-sh`. That is, when you quit `kubed-sh` all
      pre-defined environment variables are reset to their default values.

## Examples

Let's see some of the commands and env vars in action.

For example, let's say you want to launch a [simple app server in Python](https://github.com/mhausenblas/kubed-sh/blob/master/tc/python/testlr_3.py).
This app server uses port `8080` and `kubed-sh` by default exposes port `80`.
So, in order to launch it and being able to connect it, we have to tell `kubed-sh` to use the right port: `SERVICE_PORT` to the rescue:

```
[user@example.eu-west-1.eksctl.io::demo]$ SERVICE_PORT=8080
[user@example.eu-west-1.eksctl.io::demo]$ python tc/python/testlr-3.py &
[user@example.eu-west-1.eksctl.io::demo]$ ps
DPID                          SOURCE       URL
kubed-sh-1582802881222866000  testlr-3.py  testlr-3
[user@example.eu-west-1.eksctl.io::demo]$ curl testlr-3:8080
Hello from Python 3
[user@example.eu-west-1.eksctl.io::demo]$ kill kubed-sh-1582802881222866000
[user@example.eu-west-1.eksctl.io::demo]$ ps
DPID                          SOURCE       URL

[user@example.eu-west-1.eksctl.io::demo]$ 
```

Now let's go one step further and launch the same app server with a custom
image and URL. By default, `kubed-sh` would use `python:3.6-alpine3.7`
to launch a Python script, and the URL you see in `ps` would be derived from 
the script name (use the `env` command to list the current settings). Let's overwrite both:

```
[user@example.eu-west-1.eksctl.io::demo]$ PYTHON_IMAGE=python:3.6-alpine3.10
[user@example.eu-west-1.eksctl.io::demo]$ SERVICE_NAME=myappserver
[user@example.eu-west-1.eksctl.io::demo]$ python tc/python/testlr-3.py &
[user@example.eu-west-1.eksctl.io::demo]$ ps
DPID                          SOURCE       URL
kubed-sh-1582803384687189000  testlr-3.py  myappserver
[user@example.eu-west-1.eksctl.io::demo]$ curl myappserver:8080
Hello from Python 3
```

Sometimes, you want to understand what's going on under the hood, be it for
learning Kubernetes or simply troubleshooting. Use the `debug` and `literally` command:

```
[user@example.eu-west-1.eksctl.io::demo]$ debug
DEBUG mode is now on.
[user@example.eu-west-1.eksctl.io::demo]$ ps
/usr/local/bin/kubectl config current-context
in context user@example.eu-west-1.eksctl.io::demo
DPID                          SOURCE       URL
kubed-sh-1582803384687189000  testlr-3.py  myappserver
[user@example.eu-west-1.eksctl.io::demo]$ `get po,svc,deploy
/usr/local/bin/kubectl get po,svc,deploy
NAME                                                READY   STATUS      RESTARTS   AGE
pod/curljump                                        1/1     Running     0          169m
pod/kubed-sh-1582803384687189000-5d59f6bd99-kqz2c   1/1     Running     0          2m52s

NAME                  TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/myappserver   ClusterIP   10.100.96.119   <none>        8080/TCP   2m49s

NAME                                                 READY   UP-TO-DATE   AVAILABLE   AGE
deployment.extensions/kubed-sh-1582803384687189000   1/1     1            1           2m52s
```

!!! note
    Above, we didn't use `literally` directly (the long form of the command) but
    we did use \` which is its short form. Faster to type ;)

One more: on start-up `kubed-sh` auto-discovers any available `kubectl` 
[plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) and 
makes them available via the `plugin` command. This command on its own will 
list the discovered plugins and you can select one via autocompletion (TAB key):

```
[user@example.eu-west-1.eksctl.io::demo]$ plugin
view_utilization  bindrole          kboom             access_matrix     fleet             krew              kubesec_scan      rbac_view
[user@example.eu-west-1.eksctl.io::demo$ plugin fleet
CLUSTER                                               VERSION            NODES NAMESPACES PROVIDER API
scs.eu-west-1.eksctl.io                               v1.14.9-eks-c0eccc 3/3   10         AWS      https://123456789012345678901234567890AA.sk1.eu-west-1.eks.amazonaws.com
mngbase.us-west-2.eksctl.io                           v1.14.9-eks-c0eccc 2/2   6          AWS      https://123456789012345678901234567890AB.gr7.us-west-2.eks.amazonaws.com
```

!!! tip
    In above example, I executed the [fleet](https://github.com/kubectl-plus/kcf) 
    plugin, listing cluster infos. You can find more [krew](https://krew.dev) 
    plugins via the [krew index](https://github.com/kubernetes-sigs/krew-index/blob/master/plugins.md).


With these basic usage examples covered, let's have a look at some more advanced
configuration options.

## Run-time configuration

To influence the runtime behavior of `kubed-sh` you can use environment variables.

!!! note
      You can to define the environment variable in the parent shell (such as bash),
      that is, the shell you're launching `kubed-sh` itself from.

The currently supported `kubed-sh` environment variables are:

| environment variable| default | set to … |
| -------------------:| ------- | ------- |
| `KUBEDSH_GC`        | -       | define garbage collection strategy | 
| `KUBEDSH_DEBUG`     | `false` | print detailed debug messages |
| `KUBEDSH_PREPULL` | `false` | pre-pull images via `DaemonSet` |
| `KUBECTL_BINARY`    | use `which kubectl` | use this binary for API server communication |

### Speed up time-to-first-launch

If you want to speed up the time-to-first-launch, set `KUBEDSH_PREPULL=true` 
and `kubed-sh` will create a `DaemonSet` that pre-pulls the container images of
all supported languages to be ready for use. 

!!! warning
      Given the nature of a `DaemonSet`, you can use this is in all Kubernetes
      environments that explicitly allow for node access. For example, you can
      not use this feature in EKS on Fargate or OpenShift Online, where nodes 
      as such are not visible or accessible.

### Garbage collection

To influence the way `kubed-sh` performs garbage collection on exit use 
the following values for `KUBEDSH_GC`:

- `JUMP_POD` … on exit, delete the jump pod
- `ALL_PODS` … on exit, delete all pods stemming from terminating dprocs
- `ALL_DEPLOYS` …  on exit, delete all deployments and pods stemming from long-running dprocs
- `ALL_SVCS` … on exit, delete all services stemming from long-running dprocs

You can also combine the values, for example, setting `KUBEDSH_GC=JUMP_POD,ALL_DEPLOYS` 
would delete the jump pod and all deployments and pods from long-running dprocs.

## Modes

You can use `kubed-sh` either interactively or in script mode. In script mode, 
you provide `kubed-sh` a script file to interpret or make it executable—for 
example using `chmod 755 thescript` along with a [hashbang](https://en.wikipedia.org/wiki/Shebang_(Unix)) header.
The following example illustrates using `kubed-sh` in script mode.

Imagine you have a script file called `test.kbdsh` with the following content:

```
#!/usr/bin/env kubed-sh
cx user@somecluster.eu-west-1.eksctl.io
ns demo
# This line is a comment that will be ignored
node ../thescript.js &
ps
```

Above script would launch `thescript.js` in the `user@somecluster.eu-west-1.eksctl.io`
cluster, in the `demo` namespace and print out the status via `ps`.

Then, you can make it executable and execute it like so:

```
$ chmod 755 test.kbdsh
$ ./test.kbdsh
```

Alternatively you can provide a script via `stdin`:

```
$ cat tc/script.kbdsh | kubed-sh
```

… or as a command line argument:

```
$ kubed-sh tc/script.kbdsh
```

As an aside, all three above shown ways to launch a script are in fact equivalent.
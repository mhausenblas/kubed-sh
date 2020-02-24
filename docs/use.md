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

## Run-time configuration

`kubed-sh` understands the following environment variables. You have to define
them in the parent shell (such as bash) to influence the runtime behavior of 
`kubed-sh`:

| environment variable| default | set to … |
| -------------------:| ------- | ------- |
| `KUBEDSH_DEBUG`     | `false` | print detailed debug messages |
| `KUBEDSH_PREPULL` | `false` | pre-pull images via `DaemonSet` |
| `KUBECTL_BINARY`    | use `which kubectl` | use this binary for API server communication |

!!! tip
      If you want to speed up the time-to-first-launch, set 
      `KUBEDSH_PREPULL=true` and `kubed-sh` will create a `DaemonSet`, causing
      all supported languages to be ready for use. Given the nature of a `DaemonSet`
      this is available in all Kubernetes environments that explicitly allow
      for node access. For example, you can not use this feature in 
      EKS on Fargate or OpenShift Online, where nodes as such are not visible
      or accessible.


## Modes

You can use `kubed-sh` either interactively or in script mode. In script mode, 
you provide `kubed-sh` a script file to interpret or make it executable—for 
example using `chmod 755 thescript` along with a [hashbang](https://en.wikipedia.org/wiki/Shebang_(Unix)) header.
The following example illustrates using `kubed-sh` in script mode.

Imagine you have a script file called `test.kbdsh` with the following content:

```
#!/usr/bin/env kubed-sh
cx user@somecluster.eu-west-1.eksctl.io
# This line is a comment that will be ignored
node ../thescript.js &
ps
```

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
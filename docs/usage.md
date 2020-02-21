Once you've `kubed-sh` installed, launch it and you should find yourself in an interactive shell:

```sh
$ kubed-sh
Note: It seems you're running kubed-sh in a non-Linux environment (detected: darwin),
so make sure the binaries you launch are Linux binaries in ELF format.

Detected Kubernetes client in version v1.17.2 and server in version v1.14.9-eks-c0eccc

[user@somecluster.eu-west-1.eksctl.io::demo]$
```

Above, you notice that on start-up `kubed-sh` will tell you which client and server version of Kubernetes it has detected and at any point in time you are able to tell in which context (`user@somecluster.eu-west-1.eksctl.io` here) and namespace (`demo` here) you're operating.

### Built-in commands

Supported built-in commands (see also `help`) are as follows:

```
cat (local):
                        output content of file to terminal
cd (local):
                        change working directory
curl (cluster):
                        execute a curl operation in the cluster
contexts (local):
                        list available Kubernetes contexts (cluster, namespace, user tuples)
echo (local):
                        print a value or environment variable
env (local):
                        list all environment variables currently defined
exit (local):
                        leave shell
help (local):
                        list built-in commands; use help command for more details
kill (cluster):
                        stop a distributed process
literally (local):
                        execute what follows as a kubectl command
                        note that you can also prefix a line with ` to achieve the same
ls (local):
                        lists content of directory
ps (cluster):
                        list all distributed (long-running) processes in current context
pwd (local):
                        print current working directory
sleep (local):
                        sleep for specified time interval (NOP)
use (local):
                        select a certain context to work with
```

### Modes

You can use `kubed-sh` either interactively or in script mode. In script mode, you provide `kubed-sh` a script file to interpret or make it executable (for example using `chmod 755 thescript` along with a [hashbang](https://en.wikipedia.org/wiki/Shebang_(Unix)) header). The following example illustrates using `kubed-sh` in script mode:

Imagine you have a script file called `test.kbdsh` with the following content:

```
#!/usr/bin/env kubed-sh
use user@somecluster.eu-west-1.eksctl.io
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

Note that all three ways shown above are equivalent.


### Environments

`kubed-sh` supports environments and within it variables—akin to what your local shell (bash, zsh, fish) does. There are some pre-defined environment variables which influence the creation of the cluster processes you create by either specifying a binary or interpreter and script:

- `BINARY_IMAGE` (default: `alpine:3.7`) … used for executing binaries
- `NODE_IMAGE` (default: `node:12-alpine`) … used for executing Node.js scripts
- `PYTHON_IMAGE` (default: `python:3.6-alpine3.7`) … used for executing Python scripts
- `RUBY_IMAGE` (default: `ruby:2.5-alpine3.7`) … used for executing Ruby scripts
- `SERVICE_PORT` (default: `80`) … used to expose long-running processes within the cluster
- `SERVICE_NAME` (default: `""`) … used to overwrite the URL for long-running processes within the cluster
- `HOTRELOAD` (default: `false`) … used for enabling a watch on local files to trigger automatic updates on modification (EXPERIMENTAL)

!!! tip
      You can overwrite at any time any of the above environment variables to change the runtime behavior of the cluster processes you create. All changes are valid for the runtime of `kubed-sh`. That is, when you quit `kubed-sh` all pre-defined environment variables are reset to their default values.

Useful for scripting and advanced users: with the 0.5 release there are now four sub-commands to `env` (which itself simply lists the defined variables):

- `env list` … list all defined environments in current context
- `env create ENVNAME` … create a new environment called `ENVNAME`
- `env select ENVNAME` … make environment called `ENVNAME` the active one
- `env delete ENVNAME` … delete environment called `ENVNAME`

If no environment is selected, you are operating in the global environment.
See also the [design](http://kubed.sh/design) as well as [Issue #6](https://github.com/mhausenblas/kubed-sh/issues/6) for what environments are and how to work with them. Note that when you do an `env delete ENVNAME`, this environment is reaped and goes back into the global.

### Configuration

The following environment variables, defined in the parent shell (such as bash), influence the runtime behavior of `kubed-sh`. On start-up, `kubed-sh` evaluates these environment variables and enables or changes its behavior:

| env var             | default | set for |
| -------------------:| ------- | ------- |
| `KUBEDSH_DEBUG`     | `false` | print detailed messages for debug purposes |
| `KUBEDSH_NOPREPULL` | `false` | disable image pre-pull |
| `KUBECTL_BINARY`    | `which kubectl` is used to determine the binary | if set, rather than using auto-discovery, use this binary for `kubectl` commands |

!!! tip
      If you are in an environment such as EKS on Fargate or OpenShift Online where you can't create a DaemonSet, then simply launch `kubed-sh` with `$ KUBEDSH_NOPREPULL=true kubed-sh`

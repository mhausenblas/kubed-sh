[![GitHub release](https://github.com/mhausenblas/kubed-sh/workflows/release/badge.svg)](https://github.com/mhausenblas/kubed-sh/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/mhausenblas/kubed-sh)](https://goreportcard.com/report/github.com/mhausenblas/kubed-sh)

Welcome to `kubed-sh`, the Kubernetes distributed shell for the casual cluster user. In a nutshell, `kubed-sh` lets you execute a program in a Kubernetes cluster  without having to create a container image or learn new concepts. 

For example, let's say you have a Node.js script called [test.js](https://raw.githubusercontent.com/mhausenblas/kubed-sh/master/tc/node/test.js) 
and you want to launch it as a containerized app in your Kubernetes cluster, 
here's what you'd need to do in `kubed-sh`:

```
[kind::default]$ node test.js &
[kind::default]$ ps
DPID                          SOURCE      URL
kubed-sh-1517679562543558000  test.js     test
```

Does this look familiar to what you do in your "local" shell? That's the point of `kubed-sh`: it allows you to use Kubernetes without needing to learn anything new.

[Try it out](https://www.katacoda.com/mhausenblas/scenarios/kubed-sh_101) for free in your browser and/or see it in action:

| [![Introducing kubed-sh](img/introducing-kubed-sh.png)](https://www.youtube.com/watch?v=gqi1-XLiq-o) | [![kubed-sh hot-reload feature demo](img/hotreload.png)](https://www.useloom.com/share/441a97fd48ae46da8d786194f93968f6) |
|:--------------------------------:|:------------------------------------------:|
| *Introducing kubed-sh (5 min)*   | *kubed-sh hot-reload feature demo (3 min)* |

In addition to launching (Linux ELF) binaries directly, the following interpreted environments are supported:

- When you enter `node script.js`, `kubed-sh` launches a Node.js (default v12) container, copies the script into it and starts it.
- When you enter `python script.py`, `kubed-sh` launches a Python (default v3.6) container, copies the script into it and starts it.
- When you enter `ruby script.rb`, `kubed-sh` launches a Ruby (default v2.5) container, copies the script into it and starts it.

!!! tip
    Since `kubed-sh` is a proper shell environment. This means you can expect
    features such as auto-complete of built-in commands, history operations (`CTRL+R`), or clearing the screen (`CTRL+L`) to work as per usual.

Want to give it a try? Go ahead and [install](install) it now!
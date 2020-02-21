If you have access to a [Kubernetes](https://kubernetes.io/) cluster and you have [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed, you're good to go. You might want to consider using `kubed-sh`, for example for:

- **Prototyping**—Let's say you quickly want to try out a Python script or, in the context of microservices, see how a Go program and a Node.js script play together.
- **Developing**—Imagine you're developing a program in Ruby and want to launch it in a Kubernetes cluster, without having to build an image and pushing it to a registry. In this case, the experimental hot-reload feature (using `HOTRELOAD=true`) is useful for you. Whenever you save the file locally, it gets updated in the Kubernetes cluster, if hot-reload is enabled.
- **Learning Kubernetes**—You're new to Kubernetes and want to learn how to interact with it? Tip: if you issue the `debug` command you can see which `kubectl` commands `kubed-sh` launches in the background.

Also, you may be interested in [my motivation](../motivation) for writing `kubed-sh`?
To use `kubed-sh` you must meet two prerequisites:

1. `kubectl` must be [installed](https://kubernetes.io/docs/tasks/tools/install-kubectl/), tested with client version up to and included `1.17`.
1. Access to a Kubernetes cluster must be configured, tested with up to `1.14`. 

To verify your setup, you can use the following two steps:

- If you execute `ls ~/.kube/config > /dev/null && echo $?` and you see a `0` as a result, you're good, and further
- If you execute `kubectl config get-contexts | wc -l` and see a number greater than `0`, then that's super dope.

Now, download the [latest binary](https://github.com/mhausenblas/kubed-sh/releases/latest) for Linux and macOS.

For example, to install `kubed-sh` from binary on macOS you could do the following:

```sh
curl -L https://github.com/mhausenblas/kubed-sh/releases/latest/download/kubed-sh_darwin_amd64.tar.gz \
    -o kubed-sh.tar.gz && \
    tar xvzf kubed-sh.tar.gz kubed-sh && \
    mv kubed-sh /usr/local/bin && \
    rm kubed-sh*
```
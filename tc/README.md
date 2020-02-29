# Test Cases

For end-to-end tests see the [e2e](e2e/) directory.

## Observe

```
$ watch kubectl get po,svc,deploy
```

## Clean up

All processes:

```
$ kubectl delete deploy,svc,po --selector=gen=kubed-sh
```

Pre-pull only:

```
$ kubectl delete ds --selector=gen=kubed-sh,scope=pre-flight
```

## Hot-reload

```
[minikube::default] HOTRELOAD=true

---

$ cd tc
$ vi test.js
```

## Linux binary (ELF format)

In `$GOPATH/src/github.com/mhausenblas/kubed-sh/tc`:

```
$ GOOS=linux go build
$ ls -al tc
-rwxr-xr-x@ 1 mhausenblas  staff  1864063 14 Jan 06:25 tc
```

## Node.js

See `test.js` and `another.js` both long-running.

## Python

See `test.py`, one-shot.

## Ruby

See `test.rb`, one-shot.


## Katacoda

In the [kubed-sh 101](https://www.katacoda.com/mhausenblas/scenarios/kubed-sh_101) scenario, do the following:

```
curl -L https://github.com/mhausenblas/kubed-sh/releases/latest/download/kubed-sh_linux_amd64.tar.gz \
    -o kubed-sh.tar.gz && \
    tar xvzf kubed-sh.tar.gz kubed-sh && \
    mv kubed-sh /usr/local/bin && \
    rm kubed-sh*

kubectl config set-context katacoda --cluster=kubernetes --user=kubernetes-admin && kubectl config use-context katacoda
```

One last step which might be necessary: configure `kubectl` to talk to the API server:

```
kubectl config set-cluster kubernetes --server=https://172.17.0.46:6443
```

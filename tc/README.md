# Test Cases

## Observe

```
$ watch kubectl get po,svc,deploy
```

## Clean up

```
$ kubectl delete deploy,svc,po -l=gen=kubed-sh
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

# Test Cases

## Clean up

```
$ kubectl delete po,deploy,svc -l=gen=kubed-sh
```


## Linux binary (ELF format)

In `$GOPATH/src/github.com/mhausenblas/kubed-sh/tc`:

```
$ GOOS=linux go build
$ ls -al tc
-rwxr-xr-x@ 1 mhausenblas  staff  1864063 14 Jan 06:25 tc
```

## Node.js

TBD

## Python

TBD

## Ruby

TBD

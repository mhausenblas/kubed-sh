# Using `kubed-sh` containerized

To launch `kubed-sh` containerized, as a Kubernetes deployment, do:

```
$ kubectl run kubed-sh --image=quay.io/mhausenblas/kubed-sh:0.5.2
$ kubectl port-forward deploy/kubed-sh 8888:8080
```

Now you can access it via `localhost:8888` on you local machine.
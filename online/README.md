# Using `kubed-sh` containerized

To launch `kubed-sh` containerized, as a Kubernetes deployment, do:

```shell
$ kubectl create ns kd
$ ./launch.sh kd
```

Now you can access `kubed-sh` in your favorite browser at `http://localhost:8888` on you local machine.

When you're done, simply delete the namespace `kd` and with it all its resources will be removed:

```shell
$ kubectl delete ns kd
```
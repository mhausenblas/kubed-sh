package main

const (
	longRunningTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: DPROC_NAME
spec:
  replicas: 1
  selector:
    matchLabels:
      app: DPROC_NAME
  template:
    metadata:
      labels:
        app: DPROC_NAME
    spec:
      containers:
        - image: DPROC_IMAGE
          name: main
          command:
            - "sh"
            - "-c"
            - "sleep 10000"
`
	terminatingTemplate = `apiVersion: v1
kind: Pod
metadata:
  name: DPROC_NAME
  labels:
    app: DPROC_NAME
spec:
  restartPolicy: Never
  containers:
    - image: DPROC_IMAGE
      name: main
      command:
        - "sh"
        - "-c"
        - "sleep 20"
`

	prePullImgDS = `apiVersion: APIVERSION
kind: DaemonSet
metadata:
  name: PREPULLID
  annotations:
    source: "https://gist.github.com/itaysk/7bc3e56d69c4d72a549286d98fd557dd"
  labels:
    gen: kubed-sh
    scope: pre-flight
spec:
  selector:
    matchLabels:
      name: prepull
  template:
    metadata:
      labels:
        name: prepull
    spec:
      initContainers:
      - name: prepull
        image: docker
        command: ["docker", "pull", "IMG"]
        volumeMounts:
        - name: docker
          mountPath: /var/run
      volumes:
      - name: docker
        hostPath:
          path: /var/run
      containers:
      - name: pause
        image: gcr.io/google_containers/pause
`
)

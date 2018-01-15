# Motivation

While Kubernetes is a very powerful, flexible, and extensible environment to run applications, it can also be quite overwhelming, especially for casual users, effectively representing a barrier to entry.

The following issues motivated me to write `kubed-sh`:

## Container images

In order to launch an application, you first need a container image that Kubernetes in turn can run, then. There are multiple options to create a container image, from using a [CI/CD pipeline](https://fabric8.io/guide/cdelivery.html) to specialized build processes such as [S2I](https://docs.openshift.org/latest/architecture/core_concepts/builds_and_image_streams.html#source-build) to local builds and manual pushes. One can also use off-the-shelf images, launch them and then use tools like [Telepresence](https://www.telepresence.io/) or [ksync](https://vapor-ware.github.io/ksync/) to transfer the application code into the running container. In any case, one has to deal with images, directly.

## Primitives

Kubernetes introduces a number of [primitives](https://kubernetes.io/docs/concepts/) such as pods, deployments, and services, forming the building blocks for the applications you deploy and run. Now, without having at least a basic familiarity with the workloads, networking, and storage concepts it's hard to figure when to use what and how to combine things.

Some platforms, such as OpenShift, building on top of Kubernetes help here by hiding certain complex primitives and making it easier for the casual user, however oftentimes certain resource types, such as pods or services still are present and the user needs to be knowledgeable around them, nevertheless.

When coming from single-machine operating system such as Linux, you're already familiar with certain primitives such as processes, files, and job management in the shell. Put in other words, in order to execute a, say, Python script, you only need to know how to type `python thescript.py`, without studying Linux kernel structures first.

## Interactivity

For many folks, the go-to tool for interacting with a Kubernetes cluster on the command line to date is still `kubectl`. While there exists an [array of tools that make your life easier](https://abhishek-tiwari.com/10-open-source-tools-for-highly-effective-kubernetes-sre-and-ops-teams/) when working on the CLI, not many are interactive, and/or offer a fully functional shell environment.

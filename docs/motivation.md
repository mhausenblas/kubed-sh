While Kubernetes is a very powerful, flexible, and extensible environment to run applications, it can also be quite overwhelming, especially for casual users. This effectively means there's a barrier to entry for developers.

Three issues motivated me to write `kubed-sh` in the first place: 1. the number of Kubernetes primitives one needs to learn in order to use it, 2. lack of 
interactivity of existing tools, and 3. the need to build/push/pull container images.

You can read more about this topic in [As We May Kube](https://itnext.io/as-we-may-kube-293b30c0a365) or have a quick look here:


## Learning curve: primitives

Kubernetes introduces a number of [primitives](https://kubernetes.io/docs/concepts/) such as pods, deployments, and services. These form the building blocks for the applications you deploy and run. Now, without having at least a basic familiarity with the workloads, networking, and storage concepts it's hard to figure when to use what and how to combine things.

Some platforms, such as OpenShift, help here by hiding certain complex primitives and making it easier for the casual user, however oftentimes certain resource types, such as pods or services still are present and the user needs to be knowledgeable around them, nevertheless.

When coming from single-machine operating system such as Linux, you're already familiar with certain primitives such as processes, files, and job management in the shell. Put in other words, in order to execute a, say, Python script, you only need to know how to type `python thescript.py`, without studying Linux kernel structures first.


## Lack of interactivity

For many folks the go-to tool for interacting with a Kubernetes cluster on the command line to date is still `kubectl`. While there are [number of tools available](https://abhishek-tiwari.com/10-open-source-tools-for-highly-effective-kubernetes-sre-and-ops-teams/) that  make your life easier, not many are interactive and/or offer a fully functional shell environment.


## Container image build drag

In order to launch an application, you first need a container image that Kubernetes in turn can run, then. There are multiple options to create a container image, from using a [CI/CD pipeline](https://fabric8.io/guide/cdelivery.html) to specialized build processes such as [S2I](https://docs.openshift.org/latest/architecture/core_concepts/builds_and_image_streams.html#source-build) to local builds and manual pushes. One can also use off-the-shelf images, launch them and then use tools like [Telepresence](https://www.telepresence.io/) or [ksync](https://vapor-ware.github.io/ksync/) to transfer the application code into the running container. In any case, one has to deal with images, directly and explicitly. Wouldn't it be nice if that task
is automagically done for you?





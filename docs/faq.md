!!! question "Question"
    For whom is `kubed-sh`? When to use it?

!!! quote "Answer"
    I suppose it's mainly useful in a prototyping, development, or testing phase, although for low-level interactions you might find it handy in prod environments as well since it provides an interactive, context-aware version of `kubectl`. See also [use cases](../usecases/).


!!! question "Question"
    How is `kubed-sh` pronounced?

!!! quote "Answer"
    Glad you asked. Well, I pronounce it /ku:bˈdæʃ/ as in 'kube dash' ;)

!!! question "Question"
    Why another Kubernetes shell? There are already some, such as [cloudnativelabs/kube-shell](https://github.com/cloudnativelabs/kube-shell),
    [errordeveloper/kubeplay](https://github.com/errordeveloper/kubeplay), and [c-bata/kube-prompt](https://github.com/c-bata/kube-prompt). Are they not cool or what?

!!! quote "Answer"
    True, there is previous art, though these shells more or less aim at making `kubectl` interactive, exposing the commands such as `get` or `apply` to the user.
    
    In a sense `kubed-sh` is more like [technosophos/kubeshell](https://github.com/technosophos/kubeshell), trying to provide an environment a typical *nix user is comfortable with. For example, rather than providing a `create` or `apply` command to run a program, the user would simply enter the name of the executable, as she would do, for example, in the bash shell. See also the [motivation](why.md).

!!! question "Question"
    How does `kubed-sh` work? 

!!! quote "Answer"
    Good question. Essentially it's really just a glorified `kubectl` wrapper on steroids. See also the [architecture](../design/#architecture) section.
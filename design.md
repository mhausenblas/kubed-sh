# Design

The system architecture looks as follows:

![kubed-sh system architecture](img/kubed-sh-arch.png)

Components:

TBD

Overall design:

- Uses the REPL package [chzyer/readline](https://github.com/chzyer/readline) for basic shell interaction.
- Depends on `kubectl` for all cluster operations.
- A binary or script with `&` at the end causes a deployments & service, otherwise a pod is created.
- Supports environment variables to define and overwrite behavior such as exposed port, etc.

kubedsh_version := 0.4

.PHONY: build clean

build :
  # macOS binaries:
	GOOS=darwin go build -ldflags "-X main.releaseVersion=$(kubedsh_version)" -o ./kubed-sh-macos .
	# Linux binaries:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.releaseVersion=$(kubedsh_version))" -o ./kubed-sh-linux .
	# Windows binaries:
	# TBD

clean :
	@rm kubed-sh-macos
	@rm kubed-sh-linux

release_version:= 0.81

export GO111MODULE=on

.PHONY: bin
bin:
	go build -o bin/kubed-sh github.com/mhausenblas/kubed-sh

.PHONY: release
release:
	curl -sL https://git.io/goreleaser | bash -s -- --rm-dist --config .goreleaser.yml

.PHONY: publish
publish:
	git tag ${release_version}
	git push --tags
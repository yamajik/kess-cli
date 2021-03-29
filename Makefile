VERSION ?= $(shell git describe --abbrev=0 --tags master)

upgrade:
	git flow release start ${VERSION}
	git flow release finish ${VERSION}

release:
	goreleaser --rm-dist

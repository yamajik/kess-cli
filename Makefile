VERSION ?= $(shell git describe --abbrev=0 --tags master)

upgrade:
	git flow release start ${VERSION}
	git flow release finish -m "${VERSION}" ${VERSION}
	git push --all
	git push --tags

release:
	goreleaser --rm-dist

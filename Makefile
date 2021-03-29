VERSION ?= $(shell git describe --abbrev=0 --tags master)

upgrade: export GIT_MERGE_AUTOEDIT = no
upgrade:
	git flow release start ${VERSION}
	go mod tidy && go generate ./...
	git add . && git commit -m "Ver: ${UPGRADE_VERSION}" --allow-empty
	git flow release finish -m "${VERSION}" ${VERSION}
	git push --all
	git push --tags
	git checkout develop

release:
	git checkout main
	goreleaser --rm-dist
	git checkout develop

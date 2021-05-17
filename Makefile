export GO111MODULE=on

default: test

ci: test

test:
	go test ./...

build:
	go build -o kubectl-credentials_broker

depsdev:
	go get github.com/Songmu/ghch/cmd/ghch
	go get github.com/Songmu/gocredits/cmd/gocredits

prerelease: check-version
	go mod tidy
	go generate ./...
	git pull origin --tag
	ghch -w -N v${VERSION}
	gocredits . > CREDITS
	git add CHANGELOG.md CREDITS
	git commit -m "bump version to v${VERSION}"
	git tag v${VERSION}

release:
	goreleaser --rm-dist

check-version:
ifndef VERSION
	@echo 'VERSION is not set'
	@exit 1
endif
	@echo "VERSION: $$VERSION"

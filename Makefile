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

prerelease:
	git pull origin --tag
	ghch -w -N ${VERSION}
	gocredits . > CREDITS
	git add CHANGELOG.md CREDITS
	git commit -m "bump version to ${VERSION}"
	git tag ${VERSION}

release:
	goreleaser --rm-dist

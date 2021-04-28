export GO111MODULE=on

default: test

ci: test

test:
	go test ./...

build:
	go build -o kubectl-credentials_broker

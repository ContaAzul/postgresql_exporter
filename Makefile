setup:
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v
	gometalinter -i -u

build:
	go build

test:
	go test -v ./...

lint:
	gometalinter --vendor --deadline=60s ./...

ci: test lint

.DEFAULT_GOAL := build

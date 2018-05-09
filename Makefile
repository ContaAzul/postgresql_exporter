setup:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v
	gometalinter.v2 -i -u

build:
	CGO_ENABLED=0 go build -ldflags="-s -w"

test:
	go test -v ./...

lint:
	gometalinter.v2 --vendor --deadline=60s ./...

ci: test

.DEFAULT_GOAL := build

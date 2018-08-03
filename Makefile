setup:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v -vendor-only
	gometalinter.v2 --install

build:
	CGO_ENABLED=0 go build -ldflags="-s -w"

test:
	GOCACHE=off go test -v ./...

lint:
	gometalinter.v2 --vendor --deadline=60s ./...

ci: test

.DEFAULT_GOAL := build

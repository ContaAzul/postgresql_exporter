test:
	go test -v ./...

setup:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

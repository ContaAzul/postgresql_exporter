build:
	CGO_ENABLED=0 GO111MODULES=on go build -ldflags="-s -w" -o postgresql_exporter .

test:
	GOCACHE=off GO111MODULES=on go test -v ./...

lint:
	gometalinter.v2 --vendor --deadline=60s ./...

ci: test

.DEFAULT_GOAL := build

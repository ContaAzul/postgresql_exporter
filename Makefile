docker-build:
	docker build -t caninjas/postgresql_exporter .

docker-push:
	docker push caninjas/postgresql_exporter

docker: docker-build docker-push

test:
	go test -v $$(go list ./... | grep -v /vendor/)

setup:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure -v

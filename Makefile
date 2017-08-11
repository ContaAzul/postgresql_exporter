docker-build:
	@docker build -t caninjas/postgresql_exporter .

docker-push:
	@docker push caninjas/postgresql_exporter

docker: docker-build docker-push

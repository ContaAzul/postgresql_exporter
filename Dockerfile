FROM golang:alpine as build
RUN apk --no-cache add ca-certificates make
WORKDIR /go/src/app
COPY . .

FROM scratch
EXPOSE 9111
WORKDIR /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/postgresql_exporter /
ENTRYPOINT ["./postgresql_exporter"]

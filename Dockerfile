FROM golang:1.9rc2-alpine3.6 AS builder
WORKDIR /go/src/github.com/ContaAzul/postgresql_exporter
ADD . .
RUN apk add -U git
RUN go get -v github.com/golang/dep/...
RUN dep ensure -v
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build --ldflags "-extldflags "-static"" -o postgresql_exporter .

FROM scratch
EXPOSE 9111
WORKDIR /
COPY --from=builder /go/src/github.com/ContaAzul/postgresql_exporter/postgresql_exporter .
ENTRYPOINT ["./postgresql_exporter"]

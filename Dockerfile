FROM scratch
EXPOSE 9111
WORKDIR /
COPY postgresql_exporter .
ENTRYPOINT ["./postgresql_exporter"]

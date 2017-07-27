# postgresql_exporter

A Prometheus exporter for some postgresql metrics.

You can add as many database connections as you like to the
`config.yml` file, and run it with:

```console
./postgresql_exporter -config=my/config.yml
```

Then you can add hostname:9111 to the prometheus scrapes config:

```yml
- job_name: 'postgresql'
  static_configs:
    - targets: ['localhost:9111']
```

And voilá, metrics should be there and you should be able to query,
graph and alert on them.

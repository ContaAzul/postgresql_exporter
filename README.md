# postgresql_exporter

A Prometheus exporter for some postgresql metrics.

## Getting Started

You can add as many database connections as you like to the
`config.yml` file, and run it with:

```console
./postgresql_exporter -config=my/config.yml
```

By default some stat views like pg_stat_statements and pg_stat_activity doesn't allow viewing queries run by other users, unless you are a database superuser. Since you probably don't want monitoring to run as a superuser, you can setup, in a AWS RDS instance, a separate monitoring user like this:

```sql
GRANT pg_monitor TO my_monitor_user;
```

Then you can add hostname:9111 to the prometheus scrapes config:

```yml
- job_name: 'postgresql'
  static_configs:
    - targets: ['localhost:9111']
```

And voilá, metrics should be there and you should be able to query,
graph and alert on them.


## Running it within Docker

```console
docker run -p 9111 -v /path/to/my/config.yml:/config.yml caninjas/postgresql_exporter
```
# SQLC

Following the tutorial from [SQLC](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html).

Run the migration

```shell
dbmate up
```

Generate the code from the SQL queries:

```shell
sqlc generate
```

See the tools in [tools.go](tools.go) and the SQL queries in [sqlc.yaml](sqlc.yaml).
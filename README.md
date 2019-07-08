# dbjumper

> postgres connection proxy

This allows an application to connect to several db replicas

### Example config

dbjumper.yaml
```
address: "127.0.0.1:6543"
instances:
    replica1:
      type: "postgres"
      connectionstring: "postgres://postgres@127.0.0.1:5432/postgres?sslmode=disable"
    replica2:
      type: "postgres"
      connectionstring: "postgres://postgres@127.0.0.1:5432/postgres?sslmode=disable"
```

```sh

./dbjumper
```

With a custom file path:

```sh
CONFIG_PATH=/app/config ./dbjumper
```

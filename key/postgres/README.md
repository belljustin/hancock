# Hancock PostgreSQL Driver

Implements the hancock `key.Storage` interface with PostgreSQL.

## Config

Typical configuration may look like the following:

```json
{
    "backend": "postgres",
    "storage": {
        "encryption": "aes",
	"key": "secret",

        "user": "hancock",
	"password": "password",

	"host": "127.0.0.1",
	"port": 5432,
        "dbname": "hancock",

	"sslmode": "require"
    }
}
```

The database password may also be passed as environment variable by setting `HANCOCK_POSTGRES_PASSWORD`. The configuration file has precendence over environment variables.

## Running migrations

Migrations are supported by [rambler](https://github.com/elwinar/rambler).
Please install rambler:
```sh
go install github.com/elwinar/rambler
```

Then edit the `rambler.dev.config` database connections settings appropriately for your setup, and, from this directory, run:
```sh
rambler -c rambler.dev.config apply --all
```

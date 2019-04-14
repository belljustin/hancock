# Hancock PostgreSQL Driver

Provides data storage for hancock keys with PostgreSQL.

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

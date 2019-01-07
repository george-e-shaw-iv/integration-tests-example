# Integration Tests Example

This project showcases how Ardan Labs handles integration tests in regards
to databases. The patterns and techniques used to write integration tests
with a database can also be applied to other services.

This repository runs a simple list daemon (`listd`) that implements a REST API
for lists and items in relation to lists. The daemon uses a postgres database
to persist data.

For information regarding the `listd` API, see the `apiary.apib` located in
`./cmd/listd/deploy/`.

## Table of Contents

- [Running](#running)
    - [Dependencies](#dependencies)
    - [Environment Variables](#environment-variables)
    - [Make Rule](#make-rule)
- [Testing](#testing)
    - [Dependencies](#dependencies-2)
    - [Make Rule](#make-rule-2)

## Running

### Dependencies 

The only dependencies to run the services in this repository are:

- `docker`
- `docker-compose`

### Environment Variables

The program looks for the following environment variables:

- `LIST_DAEMON_PORT`: The port that the list daemon listens to/serves from (Default: `3000`).
- `DB_USER`: The postgres database username that gets used within the postgres connection
string (Default: `root`).
- `DB_PASS`: The postgres database password that gets used within the postgres connection
string (Default: `root`).
- `DB_NAME`: The postgres database name that gets used within the postgres connection string
(Default: `list`).
- `DB_HOST`: The postgres database host name that gets used within the postgres connection
string (Default `db`).
- `DB_PORT`: The postgres database port that gets used within the postgres connection string
(Default: `5432`).
- `READ_TIMEOUT`: The time, in seconds, of the read timeout of any outgoing read requests made
by the internal HTTP server (Default: `5`). 
- `WRITE_TIMEOUT`: The time, in seconds, of the read timeout of any outgoing write requests made
by the internal HTTP server (Default: `10`).
- `SHUTDOWN_TIMEOUT`: The time, in seconds, of the graceful shutdown timeout of the list daemon.
This is the amount of time in between an attempted, non-forceful shutdown and the finishing of open
requests and/or the shutdown of integrated services such as the database (Default: `5`).

If the environment variable has a supplied default and none are set within the context of the host
machine, then the default will be used.
 
To set any given environment variable, simply execute the following
pattern, replacing `[ENV_NAME]` with the name of the environment variable and `[ENV_VALUE]` with the
desired value of the environment variable: `export [ENV_NAME]=[ENV_VALUE]`. To unset any set environment
variable, simply execute the following pattern, replacing `[ENV_NAME]` with the name of the environment
variable: `unset [ENV_NAME]`.

### Make Rule

To run the services simply execute the following command:

```shell
make run
```

This will stop any containers defined by the compose file if already running
and then rebuild the containers using the compose file. The list daemon (`listd`)
will be available at `localhost:3000` and the postgres instance will be available
at `localhost:5432`.

## Testing

### Dependencies

The only dependencies to test the go code in this repository are:

- `docker`
- `docker-compose`

### Make Rule

To test the go code in this repository simply execute the following command:

```shell
make test
```

This will build the containers in docker-compose.test.yml and run
`GO111MODULE=on go test ./...` against all testable go code in the repository.
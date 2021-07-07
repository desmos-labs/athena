# DJuno

[![Codecov](https://img.shields.io/codecov/c/github/desmos-labs/djuno)](https://codecov.io/gh/desmos-labs/djuno)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/desmos-labs/djuno/Tests)](https://github.com/desmos-labs/djuno/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/desmos-labs/djuno)](https://goreportcard.com/report/github.com/desmos-labs/djuno)

This project represents the [Juno](https://github.com/desmos-labs/juno) implementation for
the [Desmos blockchain](https://github.com/desmos-labs/desmos).

It extends the custom Juno behavior with custom message handlers for all the Desmos messages. This allows to store
the needed data inside a [PostgreSQL](https://www.postgresql.org/) database on top of
which [GraphQL](https://graphql.org/) APIs can then be created using [Hasura](https://hasura.io/)

## Usage
To know how to setup and run DJuno, please refer to the [docs folder](.docs).

## Testing
If you want to test the code, you can do so by running

```shell
$ make test-unit
```

**Note**: Requires [Docker](https://docker.com).

This will:
1. Create a Docker container running a PostgreSQL database.
2. Run all the tests using that database as support.

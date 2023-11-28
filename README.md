# Athena

<p align="center" width="100%">
    <img height="90" src="./.github/logo.svg" />
</p>

<p align="center" width="100%">
  <img height="20" src="https://img.shields.io/github/license/desmos-labs/athena.svg" />
   <a href="https://codecov.io/gh/desmos-labs/athena">
      <img height="20" src="https://img.shields.io/codecov/c/github/desmos-labs/athena" />
   </a>
   <a href="https://github.com/desmos-labs/athena/actions">
      <img height="20" src="https://img.shields.io/github/actions/workflow/status/desmos-labs/athena/test.yml" />
   </a>
  <a href="https://goreportcard.com/report/github.com/desmos-labs/athena">
      <img height="20" src="https://goreportcard.com/badge/github.com/desmos-labs/athena" />
   </a>
</p>

Athena is a scraping tool for the [Desmos blockchain](https://github.com/desmos-labs/desmos) that allows to store the needed data inside a [PostgreSQL](https://www.postgresql.org/) database on top of which [GraphQL](https://graphql.org/) APIs can then be created using [Hasura](https://hasura.io/).

## Usage
To know how to setup and run Athena, please refer to the [docs folder](.docs).

## Testing
If you want to test the code, you can do so by running

```shell
$ make test-unit
```

**Note**: Requires [Docker](https://docker.com).

This will:
1. Create a Docker container running a PostgreSQL database.
2. Run all the tests using that database as support.

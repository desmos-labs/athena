# Djuno
[![Build Status](https://travis-ci.com/desmos-labs/djuno.svg?branch=master)](https://travis-ci.com/desmos-labs/djuno)
[![Go Report Card](https://goreportcard.com/badge/github.com/desmos-labs/djuno)](https://goreportcard.com/report/github.com/desmos-labs/djuno)

This project represents the [Juno](https://github.com/desmos-labs/juno) implementation for the [Desmos blockchain](https://github.com/desmos-labs/desmos).

It extends the custom Juno behavior with custom message handlers for all the Desmos messagejuno. This allows to store the needed data inside a [PostgreSQL](https://www.postgresql.org/) database on top of which [GraphQL](https://graphql.org/) APIs can then be created using [Hasura](https://hasura.io/) 

## Installation
To install the binary simply run `make install`.

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

## Database
Before running the parser, you need to: 

1. Create a [PostgreSQL](https://www.postgresql.org/) database.
2. Run the SQL queries you find inside the `*.sql` files in the [schema folder](schema) inside such database to create all the necessary tables.  

## Running the parser
To parse the chain state, you need to use the following command: 

```shell
djuno parse <path/to/config.toml>

# Example
# djuno parse config.toml 
```

The configuration must be a TOML file containing the following fields:

```toml
rpc_node = "<rpc-ip/host>:<rpc-port>"
client_node = "<client-ip/host>:<client-port>"

[database]
host = "<db-host>"
port = <db-port>
name = "<db-name>"
user = "<db-user>"
password = "<db-password>"
ssl_mode = "<ssl-mode>"
```

Example of a configuration to parse the chain state from a local full-node: 

```
rpc_node = "http://localhost:26657"
client_node = "http://localhost:1317"

[database]
host = "localhost"
port = 5432
name = "desmos"
user = "<username>"
password = "<password>"
```

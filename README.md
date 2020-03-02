# Djuno

[![Build Status](https://travis-ci.com/desmos-labs/djuno.svg?branch=master)](https://travis-ci.com/desmos-labs/djuno)

Djuno is the juno implementation for Desmos Network.  
It provides handlers that manage every incoming desmos' message and saves the successful ones inside a postgreSQL database. 

## Installation
Djuno inherit the same simple configuration of [juno](https://github.com/desmos-labs/juno)

To install the binary run `make install`.

**Note**: Requires [Go 1.13+](https://golang.org/dl/)

### Working with PostgreSQL
#### Config
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

After you have installed the binary and config the `.toml` file you can run the following command:  
`djuno parse <path/to/config.yml>`  


## Configuration
The default `config.yaml` file should look like the following:

<details>

<summary>Example config.yaml file</summary>

```yaml
chain:
  bech32_prefix: desmos
  modules:
    - authz
    - fees
    - profiles
    - relationships
    - subspaces
    - posts
    - reactions
    - reports
    - contracts
    - notifications

node:
  type: remote
  config:
    rpc:
      client_name: djuno
      address: https://rpc.morpheus.desmos.network:443
      max_connections: 10

    grpc:
      address: https://grpc.morpheus.desmos.network:443
      insecure: false

parsing:
  workers: 1
  listen_new_blocks: true
  parse_old_blocks: true
  start_height: 1

database:
  name: djuno
  host: localhost
  port: 5432
  user: user
  password: password
  max_open_connections: 15
  max_idle_connections: 10

logging:
  level: debug
  format: text

contracts:
  tips:
    code_id: 11

notifications:
  firebase_credentials_file_path: /path/to/firebase-service-account.json
  firebase_project_id: firebase-project-id
  android_channel_id: general

filters:
  supported_subspace_ids: [ 5 ]
```

</details>

## `chain`
This section contains the details of the chain configuration.

| Attribute        |   Type   | Description                            | 
|:-----------------|:--------:|:---------------------------------------|
| `modules`        | `array`  | List of modules that should be enabled |
| `bech32_prefix`  | `string` | Bech32 prefix of the addresses         | 

### Supported modules
Currently we support the followings Desmos and Cosmos SDK modules:

- `authz` to parse the data related to the Cosmos SDK `x/authz` module
- `feegrant` to parse the data related to the Cosmos SDK `x/feegrants` module
- `fees` to parse the data related to the Desmos `x/fees` module
- `profiles` to parse the data related to the Desmos `x/profiles` module
- `relationships` to parse the data related to the Desmos `x/relationships` module
- `subspaces` to parse the data related to the Desmos `x/subspaces` module
- `posts` to parse the data related to the Desmos `x/posts` module
- `reactions` to parse the data related to the Desmos `x/reactions` module
- `reports` to parse the data related to the Desmos `x/reports` module
- `contracts` to parse the data related to smart contracts

## `node`
This section contains the details of the chain node to be used in order to fetch the data.
You can reference [this page](https://github.com/forbole/juno/blob/cosmos/v0.44.x/.docs/config.md#node) for more
details.

## `parsing`
This section determines how the data will be parsed. You can
reference [this page](https://github.com/forbole/juno/blob/cosmos/v0.44.x/.docs/config.md#parsing) for more details.

## `database`
This section contains all the different configuration related to the PostgreSQL database where DJuno will write the
data. You can reference [this page](https://github.com/forbole/juno/blob/cosmos/v0.44.x/.docs/config.md#database) for
more details.

## `logging`
This section allows to configure the logging details of DJuno. You can
reference [this page](https://github.com/forbole/juno/blob/cosmos/v0.44.x/.docs/config.md#logging) for more details.

## `contracts`
If the `contracts` module is enabled, you can use this section to customize some data about the smart contracts that
will be parsed.

### `tips`
This section defines the details about the tips smart contract that should be parsed

| Attribute       |   Type    | Description                                                    | 
|:----------------|:---------:|:---------------------------------------------------------------|
| `code_id`       | `integer` | On-chan code id referring the tips smart contract to be parsed |

## `notifications`
If the `notifications` module is enabled, you can use this section to define some details about how notifications will
be sent to clients.

| Attribute                         |   Type   | Description                                                                                | 
|:----------------------------------|:--------:|:-------------------------------------------------------------------------------------------|
| `firebase_credentials_file_path`  | `string` | Path to the JSON file containing the Firebase credentials                                  |
| `firebase_project_id`             | `string` | Id of the Firebase project that should be used to send the notifications                   | 
| `android_channel_id`              | `string` | Id of the notifications channel that should be used when sending out Android notifications | 

## `filters`
If present, this section contains the details about how messages will be filtered before being parsed.

| Attribute                      |   Type   | Description                                         | 
|:-------------------------------|:--------:|:----------------------------------------------------|
| `supported_subspace_ids`       | `array`  | List of subspace id for which to parse the messages |

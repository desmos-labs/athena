CREATE TABLE contract
(
    address TEXT   NOT NULL PRIMARY KEY,
    type    TEXT   NOT NULL,
    height  BIGINT NOT NULL
);

CREATE TABLE contract_config
(
    contract_address TEXT   NOT NULL REFERENCES contract (address) PRIMARY KEY,
    config           JSONB  NOT NULL,
    height           BIGINT NOT NULL
);
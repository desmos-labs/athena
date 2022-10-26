CREATE TABLE contract
(
    address TEXT   NOT NULL PRIMARY KEY,
    type    TEXT   NOT NULL,
    config  JSONB,
    height  BIGINT NOT NULL
);
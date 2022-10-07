CREATE TABLE contract
(
    address TEXT   NOT NULL PRIMARY KEY,
    type    TEXT   NOT NULL,
    height  BIGINT NOT NULL
);
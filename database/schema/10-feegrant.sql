CREATE TABLE fee_grant
(
    granter_address TEXT   NOT NULL,
    grantee_address TEXT   NOT NULL,
    spend_limit     COIN[],
    expiration_date TIMESTAMP WITHOUT TIME ZONE,
    allowance       JSONB  NOT NULL,
    height          BIGINT NOT NULL,
    CONSTRAINT unique_fee_grant UNIQUE (granter_address, grantee_address)
);
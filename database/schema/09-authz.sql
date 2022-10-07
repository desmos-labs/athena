CREATE TABLE authz_grant
(
    granter_address TEXT                        NOT NULL,
    grantee_address TEXT                        NOT NULL,
    msg_type_url    TEXT                        NOT NULL,
    "authorization" JSONB                       NOT NULL,
    expiration      TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height          BIGINT                      NOT NULL,
    CONSTRAINT unique_msg_type_authorization UNIQUE (granter_address, grantee_address, msg_type_url)
);
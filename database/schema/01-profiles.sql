CREATE TABLE profile
(
    address       TEXT                        NOT NULL UNIQUE PRIMARY KEY,
    dtag          TEXT                        NOT NULL DEFAULT '',
    nickname      TEXT                        NOT NULL DEFAULT '',
    bio           TEXT                        NOT NULL DEFAULT '',
    profile_pic   TEXT                        NOT NULL DEFAULT '',
    cover_pic     TEXT                        NOT NULL DEFAULT '',
    creation_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    height        BIGINT                      NOT NULL
);

CREATE TABLE dtag_transfer_requests
(
    sender_address   TEXT   NOT NULL REFERENCES profile (address),
    receiver_address TEXT   NOT NULL REFERENCES profile (address),
    height           BIGINT NOT NULL,
    CONSTRAINT unique_request UNIQUE (sender_address, receiver_address)
);

CREATE TABLE profile_relationship
(
    sender_address   TEXT   NOT NULL REFERENCES profile (address),
    receiver_address TEXT   NOT NULL REFERENCES profile (address),
    subspace         TEXT   NOT NULL,
    height           BIGINT NOT NULL,
    CONSTRAINT unique_relationship UNIQUE (sender_address, receiver_address, subspace)
);

CREATE TABLE user_block
(
    blocker_address      TEXT   NOT NULL REFERENCES profile (address),
    blocked_user_address TEXT   NOT NULL REFERENCES profile (address),
    reason               TEXT,
    subspace             TEXT   NOT NULL,
    height               BIGINT NOT NULL,
    CONSTRAINT unique_blockage UNIQUE (blocker_address, blocked_user_address, subspace)
);

/* --------------------------------------------------------------------------------------------------------------- */

CREATE TABLE chain_link_chain_config
(
    id   SERIAL NOT NULL PRIMARY KEY,
    name TEXT   NOT NULL,
    CONSTRAINT unique_chain_config UNIQUE (name)
);

CREATE TABLE chain_link
(
    id                   SERIAL                      NOT NULL PRIMARY KEY,
    user_address         TEXT                        NOT NULL REFERENCES profile (address),
    external_address     TEXT                        NOT NULL,
    chain_config_id BIGINT                      NOT NULL REFERENCES chain_link_chain_config (id) ON DELETE CASCADE,
    creation_time        TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height               BIGINT                      NOT NULL,
    CONSTRAINT unique_chain_link UNIQUE (user_address, external_address, height)
);

CREATE TABLE chain_link_proof
(
    id            SERIAL NOT NULL,
    chain_link_id BIGINT NOT NULL REFERENCES chain_link (id) ON DELETE CASCADE,
    public_key    JSONB  NOT NULL,
    plain_text    TEXT   NOT NULL,
    signature     TEXT   NOT NULL,
    height        BIGINT NOT NULL,
    CONSTRAINT unique_proof_for_link UNIQUE (chain_link_id, height)
);


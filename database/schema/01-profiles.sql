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
CREATE INDEX profile_dtag_index ON profile(dtag);

CREATE TABLE dtag_transfer_requests
(
    sender_address   TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    receiver_address TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    height           BIGINT NOT NULL,
    CONSTRAINT unique_request UNIQUE (sender_address, receiver_address)
);

CREATE TABLE profile_relationship
(
    sender_address   TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    receiver_address TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    subspace         TEXT   NOT NULL,
    height           BIGINT NOT NULL,
    CONSTRAINT unique_relationship UNIQUE (sender_address, receiver_address, subspace)
);

CREATE TABLE user_block
(
    blocker_address      TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
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
CREATE INDEX chain_link_chain_config_name_index ON chain_link_chain_config(name);

CREATE TABLE chain_link
(
    id               SERIAL                      NOT NULL PRIMARY KEY,
    user_address     TEXT                        NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    external_address TEXT                        NOT NULL,
    chain_config_id  BIGINT                      NOT NULL REFERENCES chain_link_chain_config (id),
    creation_time    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height           BIGINT                      NOT NULL,
    CONSTRAINT unique_chain_link UNIQUE (chain_config_id, user_address, external_address)
);
CREATE INDEX chain_link_user_address_index ON chain_link(user_address);
CREATE INDEX chain_link_external_address_index ON chain_link(external_address);
CREATE INDEX chain_link_chain_config_id_index ON chain_link(chain_config_id);

CREATE TABLE chain_link_proof
(
    id            SERIAL NOT NULL,
    chain_link_id BIGINT NOT NULL REFERENCES chain_link (id) ON DELETE CASCADE,
    public_key    JSONB  NOT NULL,
    plain_text    TEXT   NOT NULL,
    signature     TEXT   NOT NULL,
    height        BIGINT NOT NULL,
    CONSTRAINT unique_proof_for_link UNIQUE (chain_link_id)
);

/* --------------------------------------------------------------------------------------------------------------- */

CREATE TABLE application_link
(
    id            SERIAL                      NOT NULL PRIMARY KEY,
    user_address  TEXT                        NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    application   TEXT                        NOT NULL,
    username      TEXT                        NOT NULL,
    state         TEXT                        NOT NULL,
    result        JSONB,
    creation_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height        BIGINT                      NOT NULL,
    CONSTRAINT unique_application_link UNIQUE (user_address, application, username)
);
CREATE INDEX application_link_user_address_index ON application_link(user_address);
CREATE INDEX application_link_username_index ON application_link(username);
CREATE INDEX application_link_application_index ON application_link(application);

CREATE TABLE application_link_oracle_request
(
    id                  SERIAL NOT NULL PRIMARY KEY,
    application_link_id BIGINT NOT NULL REFERENCES application_link (id) ON DELETE CASCADE,
    request_id          BIGINT NOT NULL,
    script_id           BIGINT NOT NULL,
    call_data           JSONB  NOT NULL,
    client_id           TEXT   NOT NULL,
    height              BIGINT NOT NULL,
    CONSTRAINT unique_oracle_request UNIQUE (application_link_id, client_id)
);

CREATE TABLE profiles_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX profiles_params_height_index ON profiles_params (height);
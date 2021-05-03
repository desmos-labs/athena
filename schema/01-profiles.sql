CREATE TABLE profile
(
    address       TEXT                        NOT NULL UNIQUE PRIMARY KEY,
    dtag          TEXT                        NOT NULL DEFAULT '',
    moniker       TEXT                        NOT NULL DEFAULT '',
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


CREATE TABLE relationship
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

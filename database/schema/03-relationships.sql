CREATE TABLE user_relationship
(
    creator_address      TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    counterparty_address TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    subspace_id          BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    height               BIGINT NOT NULL,
    CONSTRAINT unique_relationship UNIQUE (creator_address, counterparty_address, subspace_id)
);

CREATE TABLE user_block
(
    blocker_address TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    blocked_address TEXT   NOT NULL REFERENCES profile (address) ON DELETE CASCADE,
    reason          TEXT,
    subspace_id     BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    height          BIGINT NOT NULL,
    CONSTRAINT unique_blockage UNIQUE (blocker_address, blocked_address, subspace_id)
);
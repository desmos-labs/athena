CREATE TABLE user_relationship
(
    /* Required for Hasura links */
    row_id               SERIAL NOT NULL PRIMARY KEY,

    creator_address      TEXT   NOT NULL,
    counterparty_address TEXT   NOT NULL,
    subspace_id          BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    height               BIGINT NOT NULL,
    CONSTRAINT unique_relationship UNIQUE (creator_address, counterparty_address, subspace_id)
);

CREATE TABLE user_block
(
    blocker_address TEXT,
    blocked_address TEXT,
    reason          TEXT,
    subspace_id     BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    height          BIGINT NOT NULL,
    CONSTRAINT unique_blockage UNIQUE (blocker_address, blocked_address, subspace_id)
);
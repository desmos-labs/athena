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
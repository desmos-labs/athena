CREATE TABLE relationship
(
    sender_address   TEXT NOT NULL REFERENCES profile (address),
    receiver_address TEXT NOT NULL REFERENCES profile (address),
    subspace         TEXT NOT NULL,
    UNIQUE (sender_address, receiver_address, subspace)
);

CREATE TABLE user_block
(
    blocker_address      TEXT NOT NULL REFERENCES profile (address),
    blocked_user_address TEXT NOT NULL REFERENCES profile (address),
    reason               TEXT,
    subspace             TEXT NOT NULL,
    UNIQUE (blocker_address, blocked_user_address, subspace)
);

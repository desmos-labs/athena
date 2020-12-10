CREATE TABLE relationship
(
    sender   TEXT NOT NULL REFERENCES profile (address),
    receiver TEXT NOT NULL REFERENCES profile (address),
    subspace TEXT NOT NULL,
    UNIQUE (sender, receiver, subspace)
);

CREATE TABLE user_block
(
    blocker      TEXT NOT NULL REFERENCES profile (address),
    blocked_user TEXT NOT NULL REFERENCES profile (address),
    reason       TEXT,
    subspace     TEXT NOT NULL,
    UNIQUE (blocker, blocked_user, subspace)
);

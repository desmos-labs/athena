CREATE TABLE profile
(
    address       TEXT                        NOT NULL UNIQUE PRIMARY KEY,
    dtag          TEXT                        NOT NULL DEFAULT '',
    moniker       TEXT                        NOT NULL DEFAULT '',
    bio           TEXT                        NOT NULL DEFAULT '',
    profile_pic   TEXT                        NOT NULL DEFAULT '',
    cover_pic     TEXT                        NOT NULL DEFAULT '',
    creation_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE dtag_transfer_requests
(
    sender_address   TEXT NOT NULL REFERENCES profile (address),
    receiver_address TEXT NOT NULL REFERENCES profile (address)
);

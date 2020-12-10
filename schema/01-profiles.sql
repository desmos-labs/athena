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

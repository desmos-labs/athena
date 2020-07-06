CREATE TABLE profile
(
    address       TEXT NOT NULL UNIQUE PRIMARY KEY,
    dtag          TEXT UNIQUE,
    moniker       TEXT,
    bio           TEXT,
    profile_pic   TEXT,
    cover_pic     TEXT,
    creation_date TIMESTAMP WITHOUT TIME ZONE
);

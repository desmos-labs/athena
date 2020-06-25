CREATE TABLE profile
(
    address     TEXT NOT NULL UNIQUE PRIMARY KEY,
    moniker     TEXT,
    name        TEXT,
    surname     TEXT,
    bio         TEXT,
    profile_pic TEXT,
    cover_pic   TEXT
);

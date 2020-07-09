CREATE TABLE reaction
(
    post_id       TEXT NOT NULL REFERENCES post (id),
    owner_address TEXT NOT NULL REFERENCES profile (address),
    short_code    TEXT NOT NULL,
    value         TEXT NOT NULL,
    CONSTRAINT react_unique UNIQUE (post_id, owner_address, short_code)
);

CREATE TABLE registered_reactions
(
    owner_address TEXT NOT NULL REFERENCES profile (address),
    short_code    TEXT NOT NULL,
    value         TEXT NOT NULL,
    subspace      TEXT NOT NULL,
    CONSTRAINT registered_react_unique UNIQUE (short_code, subspace)
);

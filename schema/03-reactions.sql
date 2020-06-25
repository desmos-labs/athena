CREATE TABLE reaction
(
    post_id    TEXT NOT NULL REFERENCES post (id),
    owner      TEXT NOT NULL REFERENCES profile (address),
    short_code TEXT NOT NULL,
    value      TEXT NOT NULL,
    PRIMARY KEY (post_id, owner, short_code)
);

CREATE TABLE registered_reactions
(
    owner      TEXT NOT NULL REFERENCES profile (address),
    short_code TEXT NOT NULL,
    value      TEXT NOT NULL,
    subspace   TEXT NOT NULL,
    PRIMARY KEY (short_code, subspace)
);

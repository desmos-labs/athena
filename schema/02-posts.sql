CREATE TABLE post
(
    id              TEXT                        NOT NULL UNIQUE PRIMARY KEY,
    parent_id       TEXT REFERENCES post (id),
    message         TEXT                        NOT NULL,
    created         timestamp without time zone NOT NULL,
    last_edited     timestamp without time zone NOT NULL,
    allows_comments boolean                     NOT NULL,
    subspace        TEXT                        NOT NULL,
    creator_address TEXT                        NOT NULL REFERENCES profile (address),
    hidden          BOOLEAN                     NOT NULL DEFAULT false
);

CREATE TABLE optional_data
(
    post_id TEXT NOT NULL REFERENCES post (id),
    key     TEXT NOT NULL,
    value   TEXT NOT NULL,
    PRIMARY KEY (post_id, key)
);

CREATE TABLE attachment
(
    id        SERIAL PRIMARY KEY,
    post_id   TEXT NOT NULL REFERENCES post (id),
    uri       TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    UNIQUE (post_id, uri)
);

CREATE TABLE attachment_tag
(
    attachment_id INTEGER NOT NULL REFERENCES attachment (id),
    tag_address           TEXT    NOT NULL REFERENCES profile (address),
    UNIQUE (attachment_id, tag_address)
)

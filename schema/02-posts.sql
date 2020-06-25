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
    optional_data   jsonb                       NOT NULL DEFAULT '{}'::jsonb,
    hidden          BOOLEAN                     NOT NULL DEFAULT false
);

CREATE TABLE comment
(
    parent_id  TEXT NOT NULL REFERENCES post (id),
    comment_id TEXT NOT NULL REFERENCES post (id)
);

CREATE TABLE media
(
    post_id   TEXT NOT NULL REFERENCES post (id),
    uri       TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    PRIMARY KEY (post_id, uri)
);

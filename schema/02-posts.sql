CREATE TABLE post
(
    id              TEXT                        NOT NULL UNIQUE PRIMARY KEY,
    parent_id       TEXT REFERENCES post (id),
    message         TEXT                        NOT NULL,
    created         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    last_edited     TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    allows_comments boolean                     NOT NULL,
    subspace        TEXT                        NOT NULL,
    creator_address TEXT                        NOT NULL REFERENCES profile (address),
    optional_data   jsonb                       NOT NULL DEFAULT '{}'::jsonb,
    hidden          BOOLEAN                     NOT NULL DEFAULT false
);

CREATE TABLE media
(
    post_id   TEXT NOT NULL REFERENCES post (id),
    uri       TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    CONSTRAINT unique_post_media UNIQUE (post_id, uri)
);

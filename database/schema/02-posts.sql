CREATE TABLE post
(
    id              TEXT                        NOT NULL UNIQUE PRIMARY KEY,
    parent_id       TEXT REFERENCES post (id),
    message         TEXT                        NOT NULL,
    created         timestamp without time zone NOT NULL,
    last_edited     timestamp without time zone NOT NULL,
    disable_comments boolean                     NOT NULL,
    subspace        TEXT                        NOT NULL,
    creator_address TEXT                        NOT NULL REFERENCES profile (address),
    hidden          BOOLEAN                     NOT NULL DEFAULT false,
    height          BIGINT                      NOT NULL
);
CREATE INDEX post_height_index ON post (height);

CREATE TABLE post_attribute
(
    post_id TEXT NOT NULL REFERENCES post (id) ON DELETE CASCADE,
    key     TEXT NOT NULL,
    value   TEXT NOT NULL,
    CONSTRAINT unique_entry UNIQUE (post_id, key)
);

/* ----------------------------------------------------------------------------------------------------------------- */

CREATE TABLE post_attachment
(
    id        SERIAL PRIMARY KEY,
    post_id   TEXT NOT NULL REFERENCES post (id) ON DELETE CASCADE,
    uri       TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    CONSTRAINT unique_attachment UNIQUE (post_id, uri)
);

CREATE TABLE post_attachment_tag
(
    attachment_id INTEGER NOT NULL REFERENCES post_attachment (id) ON DELETE CASCADE,
    tag_address   TEXT    NOT NULL REFERENCES profile (address),
    CONSTRAINT unique_attachment_tag UNIQUE (attachment_id, tag_address)
);

/* ----------------------------------------------------------------------------------------------------------------- */

CREATE TABLE poll
(
    id                      SERIAL PRIMARY KEY,
    post_id                 TEXT UNIQUE                 NOT NULL REFERENCES post (id) ON DELETE CASCADE,
    question                TEXT                        NOT NULL,
    end_date                TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    allows_multiple_answers BOOLEAN                     NOT NULL,
    allows_answer_edits     BOOLEAN                     NOT NULL
);

CREATE TABLE poll_answer
(
    poll_id     INTEGER NOT NULL REFERENCES poll (id) ON DELETE CASCADE,
    answer_id   TEXT    NOT NULL,
    answer_text TEXT    NOT NULL,
    CONSTRAINT unique_answer UNIQUE (poll_id, answer_id)
);

/* ----------------------------------------------------------------------------------------------------------------- */

CREATE TABLE user_poll_answer
(
    poll_id          INTEGER NOT NULL REFERENCES poll (id) ON DELETE CASCADE,
    answer           INTEGER NOT NULL,
    answerer_address TEXT    NOT NULL REFERENCES profile (address),
    height           BIGINT  NOT NULL,
    CONSTRAINT unique_user_answer UNIQUE (poll_id, answer, answerer_address)
);

/* ----------------------------------------------------------------------------------------------------------------- */

CREATE TABLE registered_reactions
(
    owner_address TEXT   NOT NULL REFERENCES profile (address),
    short_code    TEXT   NOT NULL,
    value         TEXT   NOT NULL,
    subspace      TEXT   NOT NULL,
    height        BIGINT NOT NULL,
    CONSTRAINT registered_react_unique UNIQUE (short_code, subspace)
);

/* ----------------------------------------------------------------------------------------------------------------- */

CREATE TABLE post_reaction
(
    post_id       TEXT   NOT NULL REFERENCES post (id),
    owner_address TEXT   NOT NULL REFERENCES profile (address),
    short_code    TEXT   NOT NULL,
    value         TEXT   NOT NULL,
    height        BIGINT NOT NULL,
    CONSTRAINT react_unique UNIQUE (post_id, owner_address, short_code)
);

/* ----------------------------------------------------------------------------------------------------------------- */

CREATE TABLE post_report
(
    id               SERIAL NOT NULL,
    post_id          TEXT   NOT NULL REFERENCES post (id),
    type             TEXT   NOT NULL,
    message          TEXT,
    reporter_address TEXT   NOT NULL REFERENCES profile (address),
    height           BIGINT NOT NULL,
    CONSTRAINT unique_report UNIQUE (post_id, type, message, reporter_address)
);


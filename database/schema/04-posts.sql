CREATE TABLE post
(
    /* Required for Hasura links */
    row_id              SERIAL                      NOT NULL PRIMARY KEY,

    subspace_id         BIGINT                      NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    section_row_id      BIGINT                      NOT NULL REFERENCES subspace_section (row_id) ON DELETE CASCADE,
    id                  BIGINT                      NOT NULL,
    external_id         TEXT,
    text                TEXT,
    author_address      TEXT                        NOT NULL,
    conversation_row_id BIGINT REFERENCES post (row_id),
    reply_settings      TEXT                        NOT NULL,
    creation_date       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    last_edited_date    TIMESTAMP WITHOUT TIME ZONE,
    height              BIGINT                      NOT NULL,
    CONSTRAINT unique_subspace_post UNIQUE (subspace_id, id)
);

CREATE TABLE post_hashtag
(
    /* Required for Hasura links */
    row_id      SERIAL NOT NULL PRIMARY KEY,

    post_row_id BIGINT NOT NULL REFERENCES post (row_id) ON DELETE CASCADE,
    start_index BIGINT NOT NULL,
    end_index   BIGINT NOT NUll,
    tag         TEXT   NOT NULL
);

CREATE TABLE post_mention
(
    /* Required for Hasura links */
    row_id          SERIAL NOT NULL PRIMARY KEY,

    post_row_id     BIGINT NOT NULL REFERENCES post (row_id) ON DELETE CASCADE,
    start_index     BIGINT NOT NULL,
    end_index       BIGINT NOT NUll,
    mention_address TEXT   NOT NULL
);

CREATE TABLE post_url
(
    /* Required for Hasura links */
    row_id        SERIAL NOT NULL PRIMARY KEY,

    post_row_id   BIGINT NOT NULL REFERENCES post (row_id) ON DELETE CASCADE,
    start_index   BIGINT NOT NULL,
    end_index     BIGINT NOT NUll,
    url           TEXT   NOT NULL,
    display_value TEXT
);

CREATE TABLE post_reference
(
    /* Required for Hasura links */
    row_id         SERIAL NOT NULL PRIMARY KEY,

    post_row_id    BIGINT NOT NULL REFERENCES post (row_id) ON DELETE CASCADE,
    type           TEXT   NOT NULL,
    reference_id   BIGINT NOT NULL,
    position_index BIGINT
);

CREATE TABLE post_attachment
(
    /* Required for Hasura links */
    row_id      SERIAL NOT NULL PRIMARY KEY,

    post_row_id BIGINT NOT NULL REFERENCES post (row_id) ON DELETE CASCADE,
    id          BIGINT NOT NULL,
    content     JSONB  NOT NULL,
    height      BIGINT NOT NULL,
    CONSTRAINT unique_post_attachment UNIQUE (post_row_id, id)
);

CREATE TABLE poll_answer
(
    /* Required for Hasura links */
    row_id            SERIAL   NOT NULL PRIMARY KEY,

    attachment_row_id BIGINT   NOT NULL REFERENCES post_attachment (row_id) ON DELETE CASCADE,
    answers_indexes   BIGINT[] NOT NULL,
    user_address      TEXT     NOT NULL,
    height            BIGINT   NOT NULL,
    CONSTRAINT unique_user_answer UNIQUE (attachment_row_id, user_address)
);

CREATE TABLE posts_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX posts_params_height_index ON posts_params (height);
CREATE TABLE reaction
(
    /* Required for Hasura links */
    row_id         SERIAL NOT NULL PRIMARY KEY,

    post_row_id    BIGINT NOT NULL REFERENCES post (row_id) ON DELETE CASCADE,
    id             BIGINT NOT NULL,
    value          JSONB  NOT NULL,
    author_address TEXT   NOT NULL,
    height         BIGINT NOT NULL,
    CONSTRAINT unique_post_reaction UNIQUE (post_row_id, id)
);

CREATE TABLE subspace_registered_reaction
(
    /* Required for Hasura links */
    row_id         SERIAL NOT NULL PRIMARY KEY,

    subspace_id    BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    id             BIGINT NOT NULL,
    shorthand_code TEXT   NOT NULL,
    display_value  TEXT   NOT NULL,
    height         BIGINT NOT NULL,
    CONSTRAINT unique_subspace_registered_reaction UNIQUE (subspace_id, id)
);

CREATE TABLE subspace_registered_reaction_params
(
    subspace_id BIGINT PRIMARY KEY REFERENCES subspace (id) ON DELETE CASCADE,
    enabled     BOOLEAN NOT NULL,
    height      BIGINT  NOT NULL
);

CREATE TABLE subspace_free_text_params
(
    subspace_id BIGINT PRIMARY KEY REFERENCES subspace (id) ON DELETE CASCADE,
    enabled     BOOLEAN NOT NULL,
    max_length  BIGINT  NOT NULL,
    reg_ex      TEXT,
    height      BIGINT  NOT NULL
);
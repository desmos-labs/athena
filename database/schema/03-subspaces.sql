CREATE TABLE subspace
(
    id            BIGINT                      NOT NULL PRIMARY KEY,
    name          TEXT                        NOT NULL,
    description   TEXT,
    treasury      TEXT,
    owner         TEXT                        NOT NULL,
    creator       TEXT                        NOT NULL,
    creation_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height        BIGINT                      NOT NULL
);

CREATE TABLE subspace_user_group
(
    /* Required for Hasura links */
    row_id      SERIAL NOT NULL PRIMARY KEY,

    subspace_id BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    id          BIGINT NOT NULL,
    name        TEXT   NOT NULL,
    description TEXT,
    permissions INT    NOT NULL,
    height      BIGINT NOT NULL,
    CONSTRAINT unique_subspace_user_group UNIQUE (subspace_id, id)
);

CREATE TABLE subspace_user_group_member
(
    /* Required for Hasura links */
    row_id       SERIAL NOT NULL,

    group_row_id BIGINT NOT NULL REFERENCES subspace_user_group (row_id),
    member       TEXT   NOT NULL,
    height       BIGINT NOT NULL,
    CONSTRAINT unique_subspace_group_membership UNIQUE (group_row_id, member)
);

CREATE TABLE subspace_user_permission
(
    /* Required for Hasura links */
    row_id       SERIAL NOT NULL PRIMARY KEY,

    subspace_id  BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    user_address TEXT   NOT NULL,
    permissions  INT    NOT NULL,
    height       BIGINT NOT NULL,
    CONSTRAINT unique_subspace_permission UNIQUE (subspace_id, user_address)
);
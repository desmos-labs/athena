CREATE TABLE subspace
(
    id               BIGINT                      NOT NULL PRIMARY KEY,
    name             TEXT                        NOT NULL,
    description      TEXT,
    treasury_address TEXT,
    owner_address    TEXT                        NOT NULL,
    creator_address  TEXT                        NOT NULL,
    creation_time    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height           BIGINT                      NOT NULL
);

CREATE TABLE subspace_section
(
    /* Required for Hasura links */
    row_id        SERIAL NOT NULL PRIMARY KEY,

    subspace_id   BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    id            BIGINT NOT NULL,
    parent_row_id BIGINT REFERENCES subspace_section (row_id) ON DELETE CASCADE,
    name          TEXT   NOT NULL,
    description   TEXT,
    height        BIGINT NOT NULL,
    CONSTRAINT unique_subspace_section UNIQUE (subspace_id, id)
);

CREATE TABLE subspace_user_group
(
    /* Required for Hasura links */
    row_id         SERIAL NOT NULL PRIMARY KEY,

    subspace_id    BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    section_row_id BIGINT NOT NULL REFERENCES subspace_section (row_id) ON DELETE CASCADE,
    id             BIGINT NOT NULL,
    name           TEXT   NOT NULL,
    description    TEXT,
    permissions    TEXT[] NOT NULL,
    height         BIGINT NOT NULL,
    CONSTRAINT unique_subspace_user_group UNIQUE (subspace_id, id)
);

CREATE TABLE subspace_user_group_member
(
    /* Required for Hasura links */
    row_id         SERIAL NOT NULL,

    group_row_id   BIGINT NOT NULL REFERENCES subspace_user_group (row_id) ON DELETE CASCADE,
    member_address TEXT   NOT NULL,
    height         BIGINT NOT NULL,
    CONSTRAINT unique_subspace_group_membership UNIQUE (group_row_id, member_address)
);

CREATE TABLE subspace_user_permission
(
    /* Required for Hasura links */
    row_id         SERIAL NOT NULL PRIMARY KEY,

    section_row_id BIGINT NOT NULL REFERENCES subspace_section (row_id) ON DELETE CASCADE,
    user_address   TEXT   NOT NULL,
    permissions    TEXT[] NOT NULL,
    height         BIGINT NOT NULL,
    CONSTRAINT unique_subspace_permission UNIQUE (section_row_id, user_address)
);
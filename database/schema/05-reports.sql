CREATE TABLE subspace_report_reason
(
    /* Required for Hasura links */
    row_id      SERIAL NOT NULL PRIMARY KEY,

    subspace_id BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    id          BIGINT NOT NULL,
    title       TEXT   NOT NULL,
    description TEXT,
    height      BIGINT NOT NULL,
    CONSTRAINT unique_subspace_reason UNIQUE (subspace_id, id)
);

CREATE TABLE report
(
    /* Required for Hasura links */
    row_id           SERIAL                      NOT NULL PRIMARY KEY,

    subspace_id      BIGINT                      NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    id               BIGINT                      NOT NULL,
    message          TEXT,
    reporter_address TEXT                        NOT NULL,
    target           JSONB                       NOT NULL,
    creation_date    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height           BIGINT                      NOT NULL,
    CONSTRAINT unique_subspace_report UNIQUE (subspace_id, id)
);

CREATE TABLE report_reason
(
    report_row_id BIGINT NOT NULL REFERENCES report (row_id) ON DELETE CASCADE,
    reason_row_id BIGINT NOT NULL REFERENCES subspace_report_reason (row_id) ON DELETE CASCADE,
    CONSTRAINT unique_report_reason UNIQUE (report_row_id, reason_row_id)
);

CREATE TABLE reports_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX reports_params_height_index ON reports_params (height);
CREATE TABLE application_link_score
(
    id                      SERIAL                      NOT NULL,
    application_link_row_id BIGINT                      NOT NULL REFERENCES application_link (id) ON DELETE CASCADE UNIQUE,
    details                 JSONB                       NOT NULL,
    score                   INT                         NOT NULL DEFAULT 0,
    timestamp               TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
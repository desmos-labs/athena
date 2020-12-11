CREATE TABLE report
(
    id               SERIAL NOT NULL,
    post_id          TEXT   NOT NULL REFERENCES post (id),
    type             TEXT   NOT NULL,
    message          TEXT,
    reporter_address TEXT   NOT NULL REFERENCES profile (address)
)

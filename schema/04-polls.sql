CREATE TABLE poll
(
    id                      SERIAL PRIMARY KEY,
    post_id                 TEXT UNIQUE                 NOT NULL REFERENCES post (id),
    question                TEXT                        NOT NULL,
    end_date                TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    allows_multiple_answers BOOLEAN                     NOT NULL,
    allows_answer_edits     BOOLEAN                     NOT NULL
);

CREATE TABLE poll_answer
(
    poll_id     INTEGER NOT NULL REFERENCES poll (id),
    answer_id   TEXT    NOT NULL,
    answer_text TEXT    NOT NULL,
    UNIQUE (poll_id, answer_id)
);

CREATE TABLE user_poll_answer
(
    poll_id          INTEGER NOT NULL REFERENCES poll (id),
    answer           INTEGER NOT NULL,
    answerer_address TEXT    NOT NULL REFERENCES profile (address),
    UNIQUE (poll_id, answer, answerer_address)
);

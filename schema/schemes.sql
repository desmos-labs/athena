--- COSMOS ----------------------------------------------

CREATE TABLE validator
(
    id               SERIAL PRIMARY KEY,
    address          character varying(40) NOT NULL UNIQUE,
    consensus_pubkey character varying(83) NOT NULL UNIQUE
);

CREATE TABLE pre_commit
(
    id                SERIAL PRIMARY KEY,
    validator_address character varying(40)       NOT NULL REFERENCES validator (address),
    timestamp         timestamp without time zone NOT NULL,
    voting_power      integer                     NOT NULL,
    proposer_priority integer                     NOT NULL
);

CREATE TABLE block
(
    id               SERIAL PRIMARY KEY,
    height           integer                     NOT NULL UNIQUE,
    hash             character varying(64)       NOT NULL UNIQUE,
    num_txs          integer DEFAULT 0,
    total_gas        integer DEFAULT 0,
    proposer_address character varying(40)       NOT NULL REFERENCES validator (address),
    pre_commits      integer                     NOT NULL,
    timestamp        timestamp without time zone NOT NULL
);

CREATE TABLE transaction
(
    id         SERIAL PRIMARY KEY,
    timestamp  timestamp without time zone NOT NULL,
    gas_wanted integer                              DEFAULT 0,
    gas_used   integer                              DEFAULT 0,
    height     integer                     NOT NULL REFERENCES block (height),
    txhash     character varying(64)       NOT NULL UNIQUE,
    messages   jsonb                       NOT NULL DEFAULT '[]'::jsonb,
    fee        jsonb                       NOT NULL DEFAULT '{}'::jsonb,
    signatures jsonb                       NOT NULL DEFAULT '[]'::jsonb,
    memo       character varying(256)
);

--- DESMOS ----------------------------------------------
​
CREATE TABLE "user"
(
    id      SERIAL PRIMARY KEY,
    address character varying(45) NOT NULL
);


CREATE TABLE poll
(
    id                      SERIAL PRIMARY KEY,
    question                text                     NOT NULL,
    end_date                timestamp with time zone NOT NULL,
    open                    boolean                  NOT NULL,
    allows_multiple_answers boolean                  NOT NULL,
    allows_answer_edits     boolean                  NOT NULL
);
​
CREATE TABLE poll_answer
(
    poll_id     integer NOT NULL REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    answer_id   integer NOT NULL,
    answer_text text    NOT NULL,
    UNIQUE (poll_id, answer_id)
);
​
CREATE TABLE user_poll_answer
(
    poll_id integer NOT NULL REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    answer  integer NOT NULL,
    user_id integer NOT NULL REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE (poll_id, answer, user_id)
);

CREATE TABLE post
(
    id              integer PRIMARY KEY,
    parent_id       integer                  NOT NULL,
    message         text                     NOT NULL,
    created         timestamp with time zone NOT NULL,
    last_edited     timestamp with time zone NOT NULL,
    allows_comments boolean                  NOT NULL,
    subspace        text                     NOT NULL,
    creator_id      integer                  NOT NULL REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE,
    optional_data   jsonb                    NOT NULL DEFAULT '{}'::jsonb,
    poll_id         integer REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE
);
​
CREATE TABLE reaction
(
    id       SERIAL PRIMARY KEY,
    post_id  integer NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE,
    owner_id integer NOT NULL REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE,
    value    text    NOT NULL
);
​​
CREATE TABLE media
(
    id        SERIAL PRIMARY KEY,
    post_id   integer NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE,
    uri       text    NOT NULL,
    mime_type text    NOT NULL
);

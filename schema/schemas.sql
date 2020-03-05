-- Post ----------------------------------------------

CREATE TABLE post
(
    id              integer PRIMARY KEY,
    parent_id       integer                  NOT NULL,
    message         text                     NOT NULL,
    created         timestamp with time zone NOT NULL,
    last_edited     timestamp with time zone NOT NULL,
    allows_comments boolean                  NOT NULL,
    subspace        text                     NOT NULL,
    creator         character varying(45)    NOT NULL,
    poll_id         integer REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    optional_data   jsonb                    NOT NULL DEFAULT '{}'::jsonb
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX post_pkey ON post (id int4_ops);

-- Poll ----------------------------------------------

CREATE TABLE poll
(
    id                      SERIAL PRIMARY KEY,
    question                text                     NOT NULL,
    end_date                timestamp with time zone NOT NULL,
    open                    boolean                  NOT NULL,
    allows_multiple_answers boolean                  NOT NULL,
    allows_answer_edits     boolean                  NOT NULL
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX poll_pkey ON poll (id int4_ops);


-- Poll_answer ----------------------------------------------

CREATE TABLE poll_answer
(
    id          SERIAL PRIMARY KEY,
    poll_id     integer NOT NULL REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    answer_id   integer NOT NULL,
    answer_text text    NOT NULL
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX poll_answer_pkey ON poll_answer (id int4_ops);

-- User_poll_answer ------------------------------------------------

CREATE TABLE user_poll_answer
(
    id           SERIAL PRIMARY KEY,
    poll_id      integer               NOT NULL REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    answers      integer[]             NOT NULL,
    user_address character varying(45) NOT NULL
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX user_poll_answer_pkey ON user_poll_answer (id int4_ops);


-- Reaction ----------------------------------------------

CREATE TABLE reaction
(
    id      SERIAL PRIMARY KEY,
    post_id integer               NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE,
    owner   character varying(45) NOT NULL,
    value   text                  NOT NULL
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX reaction_pkey ON reaction (id int4_ops);

-- Media ----------------------------------------------

CREATE TABLE media
(
    id        SERIAL PRIMARY KEY,
    post_id   integer NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE,
    uri       text    NOT NULL,
    mime_type text    NOT NULL
);

-- Indices -------------------------------------------------------

CREATE UNIQUE INDEX media_pkey ON media (id int4_ops);
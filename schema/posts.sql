CREATE TABLE "user"
(
    id          SERIAL PRIMARY KEY,
    address     character varying(45) UNIQUE NOT NULL,
    moniker     text,
    name        text,
    surname     text,
    bio         text,
    profile_pic text,
    cover_pic   text
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

CREATE TABLE poll_answer
(
    poll_id     integer NOT NULL REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    answer_id   integer NOT NULL,
    answer_text text    NOT NULL,
    UNIQUE (poll_id, answer_id)
);

CREATE TABLE user_poll_answer
(
    poll_id integer NOT NULL REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    answer  integer NOT NULL,
    user_id integer NOT NULL REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE,
    UNIQUE (poll_id, answer, user_id)
);

CREATE TABLE post
(
    id              text PRIMARY KEY,
    parent_id       text REFERENCES post (id),
    message         text                        NOT NULL,
    created         timestamp without time zone NOT NULL,
    last_edited     timestamp without time zone NOT NULL,
    allows_comments boolean                     NOT NULL,
    subspace        text                        NOT NULL,
    creator_id      integer                     NOT NULL REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE,
    optional_data   jsonb                       NOT NULL DEFAULT '{}'::jsonb,
    poll_id         integer REFERENCES poll (id) ON DELETE CASCADE ON UPDATE CASCADE,
    hidden          BOOLEAN                     NOT NULL DEFAULT false
);

CREATE TABLE comment
(
    parent_id  text NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE,
    comment_id text NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE reaction
(
    id         SERIAL PRIMARY KEY,
    post_id    text    NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE,
    owner_id   integer NOT NULL REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE,
    short_code text    NOT NULL,
    value      text    NOT NULL
);

CREATE TABLE registered_reactions
(
    id         SERIAL PRIMARY KEY,
    owner_id   integer NOT NULL REFERENCES "user" (id) ON DELETE CASCADE ON UPDATE CASCADE,
    short_code text    NOT NULL,
    value      text    NOT NULL,
    subspace   text    NOT NULL
);

CREATE TABLE media
(
    id        SERIAL PRIMARY KEY,
    post_id   text NOT NULL REFERENCES post (id) ON DELETE CASCADE ON UPDATE CASCADE,
    uri       text NOT NULL,
    mime_type text NOT NULL
);

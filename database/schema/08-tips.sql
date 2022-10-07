CREATE TABLE tip_user
(
    sender_address   TEXT   NOT NULL,
    receiver_address TEXT   NOT NULL,
    subspace_id      BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    amount           COIN[] DEFAULT '{}',
    height           BIGINT NOT NULL,
    CONSTRAINT unique_sender_user_tip UNIQUE (sender_address, receiver_address, height)
);

CREATE TABLE tip_post
(
    sender_address TEXT   NOT NULL,
    subspace_id    BIGINT NOT NULL REFERENCES subspace (id) ON DELETE CASCADE,
    post_row_id    BIGINT NOT NULL REFERENCES post (row_id) ON DELETE CASCADE,
    amount         COIN[] DEFAULT '{}',
    height         BIGINT NOT NULL,
    CONSTRAINT unique_sender_post_tip UNIQUE (sender_address, post_row_id, height)
);
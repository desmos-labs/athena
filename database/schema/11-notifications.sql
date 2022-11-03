CREATE TABLE notification
(
    user_address TEXT                        NOT NULL,
    type         TEXT                        NOT NULL,
    data         JSONB                       NOT NULL,
    timestamp    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT unique_user_notification UNIQUE (user_address, data)
);
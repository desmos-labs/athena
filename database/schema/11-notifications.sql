CREATE TABLE notification
(
    user_address TEXT                        NOT NULL,
    type         TEXT                        NOT NULL,
    data         JSONB                       NOT NULL,
    timestamp    TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT unique_user_notification UNIQUE (user_address, data)
);

CREATE TABLE notification_token
(
    user_address TEXT                        NOT NULL,
    device_token TEXT                        NOT NULL,
    timestamp    TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
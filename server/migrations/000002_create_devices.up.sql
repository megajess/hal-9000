CREATE TABLE devices (
    id            TEXT      PRIMARY KEY,
    name          TEXT      NOT NULL,
    api_key       TEXT      NOT NULL UNIQUE,
    current_state TEXT      NOT NULL DEFAULT 'off',
    desired_state TEXT      NOT NULL DEFAULT 'off',
    last_seen     DATETIME,
    created_at    DATETIME  NOT NULL
);
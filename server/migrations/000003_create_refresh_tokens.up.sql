CREATE TABLE refresh_tokens (
    token      TEXT      PRIMARY KEY,
    user_id    TEXT      NOT NULL,
    expires_at DATETIME  NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
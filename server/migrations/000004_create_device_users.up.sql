CREATE TABLE device_users (
    device_id  TEXT  NOT NULL,
    user_id    TEXT  NOT NULL,
    role       TEXT  NOT NULL,
    PRIMARY KEY (device_id, user_id),
    FOREIGN KEY (device_id) REFERENCES devices(id),
    FOREIGN KEY (user_id)   REFERENCES users(id)
);
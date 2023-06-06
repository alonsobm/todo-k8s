CREATE TABLE IF NOT EXISTS users
(
    user_id         bigserial NOT NULL,
    username        text      NOT NULL,
    email           text      NOT NULL,
    hashed_password text      NOT NULL,
    created_at      timestamp NOT NULL DEFAULT NOW(),
    updated_at      timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id),
    UNIQUE (username),
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS todos
(
    todo_id    bigserial NOT NULL,
    user_id    bigserial NOT NULL,
    title      text      NOT NULL,
    content    text      NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    PRIMARY KEY (todo_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE
);


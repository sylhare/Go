CREATE TABLE authors (
    id   BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    bio  text
);

CREATE TABLE books (
                       id          BIGSERIAL PRIMARY KEY,
                       title       TEXT NOT NULL,
                       author_id   BIGINT NOT NULL,
                       created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       FOREIGN KEY (author_id) REFERENCES authors (id)
);
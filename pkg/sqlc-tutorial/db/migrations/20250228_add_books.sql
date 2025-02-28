-- migrate:up
CREATE TABLE IF NOT EXISTS books (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    author_id   BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (author_id) REFERENCES authors (id)
);

-- migrate:down
DROP TABLE IF EXISTS books;

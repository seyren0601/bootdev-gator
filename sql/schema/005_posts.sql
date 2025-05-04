-- +goose Up
CREATE TABLE posts(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title text NOT NULL,
    url text UNIQUE NOT NULL,
    description text NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id uuid NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
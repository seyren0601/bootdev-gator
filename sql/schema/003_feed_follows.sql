-- +goose Up
CREATE TABLE feed_follows(
    user_id uuid,
    feed_id uuid,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
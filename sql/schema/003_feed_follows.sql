-- +goose Up
ALTER TABLE feeds DROP COLUMN user_id;

CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);


-- +goose Down
DROP TABLE feed_follows;

ALTER TABLE feeds
    ADD COLUMN user_id UUID NOT NULL,
    ADD CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

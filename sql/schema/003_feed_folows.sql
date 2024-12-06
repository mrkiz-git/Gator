-- +goose Up

CREATE TABLE feed_follows (
    id UUID PRIMARY KEY, -- Unique identifier for each follow
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Foreign key to users table
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE, -- Foreign key to feeds table
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- Timestamp for when the follow was created
    updated_at TIMESTAMP NOT NULL DEFAULT now(), -- Timestamp for when the follow was last updated
    UNIQUE (user_id, feed_id) -- Ensure no duplicate follow records
);

-- +goose Down

DROP TABLE fecded_follows;
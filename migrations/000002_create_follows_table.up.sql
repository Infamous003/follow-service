CREATE TABLE IF NOT EXISTS follows (
    follower_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, followee_id),
    CHECK (follower_id <> followee_id)
);
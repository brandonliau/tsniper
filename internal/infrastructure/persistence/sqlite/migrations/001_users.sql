CREATE TABLE IF NOT EXISTS users (
    user_id             TEXT,
    campus              TEXT,
    UNIQUE (user_id)
);

CREATE INDEX IF NOT EXISTS idx_users_userid ON users(user_id);

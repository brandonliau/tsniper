CREATE TABLE IF NOT EXISTS snipes (user_id TEXT, course_index TEXT, campus TEXT, season TEXT);

CREATE INDEX IF NOT EXISTS idx_snipes_userid ON snipes(user_id);

CREATE TABLE IF NOT EXISTS snipes (course_index TEXT, user_id TEXT, campus TEXT, season TEXT);

CREATE INDEX IF NOT EXISTS idx_snipes_userid ON snipes(user_id);

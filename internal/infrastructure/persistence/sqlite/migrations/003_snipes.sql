CREATE TABLE IF NOT EXISTS snipes (
    user_id      TEXT,
    course_index TEXT,
    campus       TEXT,
    term         TEXT,
    year         TEXT,
    UNIQUE(user_id, course_index, campus, term, year)
);

CREATE INDEX IF NOT EXISTS idx_snipes_user_id ON snipes(user_id);
CREATE INDEX IF NOT EXISTS idx_snipes_course_index ON snipes(course_index);

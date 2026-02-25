CREATE TABLE IF NOT EXISTS courses (
    course_index  TEXT,
    course_string TEXT,
    section       TEXT,
    title         TEXT,
    instructors   TEXT,
    notes         TEXT,
    meeting       TEXT,
    campus        TEXT,
    term          TEXT,
    year          TEXT,
    last_open     BIGINT,
    UNIQUE (course_index, campus, term, year)
);

CREATE INDEX IF NOT EXISTS idx_courses_course_index ON courses(course_index);

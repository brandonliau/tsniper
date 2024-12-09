CREATE TABLE IF NOT EXISTS courses (
    course_index TEXT,
    title TEXT,
    course_string TEXT,
    section TEXT,
    instructors TEXT,
    notes TEXT,
    meeting TEXT,
    campus TEXT,
    season TEXT,
    last_open REAL,
    UNIQUE (course_index, campus, season)
);

CREATE INDEX IF NOT EXISTS idx_courses_courseindex ON courses(course_index);
CREATE INDEX IF NOT EXISTS idx_courses_season ON courses(season);
CREATE INDEX IF NOT EXISTS idx_courses_campus_season ON courses(campus, season);

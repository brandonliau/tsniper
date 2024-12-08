CREATE TABLE IF NOT EXISTS NB (
    course_index TEXT,
    title TEXT,
    course_string TEXT,
    section TEXT,
    instructors TEXT,
    notes TEXT,
    meeting TEXT,
    season TEXT,
    last_open REAL,
    UNIQUE (course_index, season)
);
CREATE INDEX IF NOT EXISTS idx_NB_courseindex ON NB(course_index);
CREATE INDEX IF NOT EXISTS idx_NB_season ON NB(season);

CREATE TABLE IF NOT EXISTS NK (
    course_index TEXT,
    title TEXT,
    course_string TEXT,
    section TEXT,
    instructors TEXT,
    notes TEXT,
    meeting TEXT,
    season TEXT,
    last_open REAL,
    UNIQUE (course_index, season)
);
CREATE INDEX IF NOT EXISTS idx_NK_courseindex ON NK(course_index);
CREATE INDEX IF NOT EXISTS idx_NK_season ON NK(season);

CREATE TABLE IF NOT EXISTS CM (
    course_index TEXT,
    title TEXT,
    course_string TEXT,
    section TEXT,
    instructors TEXT,
    notes TEXT,
    meeting TEXT,
    season TEXT,
    last_open REAL,
    UNIQUE (course_index, season)
);
CREATE INDEX IF NOT EXISTS idx_CM_courseindex ON CM(course_index);
CREATE INDEX IF NOT EXISTS idx_CM_season ON CM(season);

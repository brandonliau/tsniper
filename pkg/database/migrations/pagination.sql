CREATE TABLE IF NOT EXISTS pagination (check_hash TEXT, check_data TEXT, date_time BIGINT);

CREATE INDEX IF NOT EXISTS idx_pagination_checkhash ON pagination(check_hash);

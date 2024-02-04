CREATE TABLE IF NOT EXISTS institution_working_hours (
    institution_id BIGINT REFERENCES institution(id) ON DELETE CASCADE,
    day_of_week INT, -- 0 = Sunday, 1 = Monday, ..., 6 = Saturday
    open_time TIME,
    close_time TIME
);

CREATE INDEX IF NOT EXISTS idx_institution_working_hours_institution_id ON institution_working_hours(institution_id);

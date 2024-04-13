
CREATE TABLE IF NOT EXISTS employee_schedule (
    employee_id BIGINT REFERENCES employee (id) ON DELETE CASCADE,
    day_of_week INT, -- 0 = Sunday, 1 = Monday, ..., 6 = Saturday
    start_time TIME,
    end_time TIME,
    break_start_time TIME,
    break_end_time TIME
);

CREATE INDEX IF NOT EXISTS idx_employee_id ON employee_schedule (employee_id);

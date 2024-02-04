CREATE TABLE IF NOT EXISTS employee_service (
    employee_id BIGINT REFERENCES employee (id) ON DELETE CASCADE,
    service_id BIGINT REFERENCES services (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_employee_id ON employee_service (employee_id);

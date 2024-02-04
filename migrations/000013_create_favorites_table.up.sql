CREATE TABLE IF NOT EXISTS favorites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,
    institution_id BIGINT REFERENCES institution(id),
    employee_id BIGINT REFERENCES employee(id) NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (service_id) REFERENCES services (id)
);

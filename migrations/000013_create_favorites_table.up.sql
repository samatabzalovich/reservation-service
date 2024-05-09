CREATE TABLE IF NOT EXISTS favorites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,
    institution_id BIGINT REFERENCES institution(id) on delete  cascade ,
    employee_id BIGINT REFERENCES employee(id)  on delete  cascade NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON delete cascade,
    FOREIGN KEY (service_id) REFERENCES services (id) on DELETE CASCADE
);

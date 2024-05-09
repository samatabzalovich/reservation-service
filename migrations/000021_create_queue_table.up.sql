Drop type if exists QUEUESTATUS;
Create TYPE QUEUESTATUS as ENUM ('pending', 'in-progress', 'completed', 'cancelled', 'called');
CREATE TABLE IF NOT EXISTS queue (
                                     id BIGSERIAL PRIMARY KEY,
                                     institution_id BIGINT REFERENCES institution(id) ON DELETE  CASCADE ,
                                     service_id BIGINT REFERENCES services(id) ON DELETE  CASCADE ,
                                     position INTEGER NOT NULL,
                                     status QUEUESTATUS NOT NULL DEFAULT 'pending' ,
    employee_id BIGINT REFERENCES employee(id) ON DELETE CASCADE NULL DEFAULT NULL ,
    user_id BIGINT REFERENCES users(id) ON DELETE  CASCADE ,
    version INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
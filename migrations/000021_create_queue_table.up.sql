Drop type if exists QUEUESTATUS;
Create TYPE QUEUESTATUS as ENUM ('pending', 'in-progress', 'completed', 'cancelled');
CREATE TABLE IF NOT EXISTS queue (
                                     id BIGSERIAL PRIMARY KEY,
                                     institution_id BIGINT REFERENCES institution(id),
                                     service_id BIGINT REFERENCES services(id),
                                     position INTEGER NOT NULL,
                                     status QUEUESTATUS NOT NULL DEFAULT 'pending' ,
    user_id BIGINT REFERENCES users(id)
);
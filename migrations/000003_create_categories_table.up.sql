CREATE TABLE IF NOT EXISTS category (
                                        id BIGSERIAL PRIMARY KEY,
                                        name VARCHAR(50) NOT NULL,
                                        description VARCHAR(255) NOT NULL,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE IF NOT EXISTS services (
    id BIGSERIAL PRIMARY KEY,
    institution_id BIGINT REFERENCES institution(id) ON DELETE CASCADE,
    name varchar(50) NOT NULL,
    description TEXT,
    duration INTERVAL NOT NULL,
    price INTEGER NOT NULL,
    photo_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


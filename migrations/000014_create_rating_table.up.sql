CREATE TABLE IF NOT EXISTS rating (
    id BIGSERIAL PRIMARY KEY,
    appointment_id BIGINT REFERENCES appointments(id),
    institution_id BIGINT REFERENCES institution(id),
    user_id BIGINT REFERENCES users(id),
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

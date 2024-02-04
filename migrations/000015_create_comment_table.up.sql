CREATE TABLE IF NOT EXISTS comment (
    id BIGSERIAL PRIMARY KEY,
    institution_id BIGINT REFERENCES institution(id),
    user_id BIGINT REFERENCES users(id),
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

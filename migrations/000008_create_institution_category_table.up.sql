CREATE TABLE IF NOT EXISTS institution_category (
    inst_id BIGINT REFERENCES institution (id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES category (id) ON DELETE CASCADE
);


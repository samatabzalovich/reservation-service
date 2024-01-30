CREATE TABLE IF NOT EXISTS institution (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    website TEXT NOT NULL,
    owner_id BIGINT REFERENCES users(id) NOT NULL,
    latitude VARCHAR(255) NOT NULL,
    longitude VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(12) NOT NULL,
    country VARCHAR(50) NOT NULL,
    city VARCHAR(50) NOT NULL,
    category_id INTEGER NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES category(id)
);

-- INSERT INTO institution (name, description, website, owner_id, latitude,
--                          longitude, address, phone, country, city, category_id)
-- VALUES ('Instituto de Computação',
--         'Instituto de Computação da Universidade Federal de Alagoas', 'https://www.ic.ufal.br/',
--         1, '-9.555067', '-35.779773',
--         'Av. Lourival Melo Mota, s/n - Cidade Universitária, ' ||
--         'Maceió - AL, 57072-900', '+77375452529',
--         'Brasil', 'Maceió', 3);

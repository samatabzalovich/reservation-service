CREATE TABLE IF NOT EXISTS category (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO category (name, description) VALUES ('Барбершоп', 'Мужская парикмахерская');
INSERT INTO category (name, description) VALUES ('Парикмахерская', 'Женская парикмахерская');
INSERT INTO category (name, description) VALUES ('Салон красоты', 'Салон красоты');
INSERT INTO category (name, description) VALUES ('Стоматолог', 'Стоматология');


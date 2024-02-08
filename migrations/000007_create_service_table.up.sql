CREATE TABLE IF NOT EXISTS services (
    id BIGSERIAL PRIMARY KEY,
    institution_id BIGINT REFERENCES institution(id) ON DELETE CASCADE,
    name varchar(50) NOT NULL,
    description TEXT,
    duration INTERVAL NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    photo_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- INSERT INTO services (institution_id, name, description, duration, price, photo_url) VALUES
--                                                                                          (1, 'Corte de cabelo', 'Corte de cabelo masculino', '00:30:00', 20.00, 'https://www.google.com.br'),
--                                                                                          (1, 'Corte de cabelo', 'Corte de cabelo feminino', '00:30:00', 30.00, 'https://www.google.com.br');

CREATE TABLE IF NOT EXISTS city (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    state VARCHAR(50) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO city (name, state) VALUES ('Астана', 'Акмолинская область');
INSERT INTO city (name, state) VALUES ('Актобе', 'Актюбинская область');
INSERT INTO city (name, state) VALUES ('Алматы', 'Алматинская область');
INSERT INTO city (name, state) VALUES ('Атырау', 'Атырауская область');
INSERT INTO city (name, state) VALUES ('Караганда', 'Карагандинская область');
INSERT INTO city (name, state) VALUES ('Кокшетау', 'Акмолинская область');
INSERT INTO city (name, state) VALUES ('Костанай', 'Костанайская область');
INSERT INTO city (name, state) VALUES ('Кызылорда', 'Кызылординская область');
INSERT INTO city (name, state) VALUES ('Павлодар', 'Павлодарская область');
INSERT INTO city (name, state) VALUES ('Петропавловск', 'Северо-Казахстанская область');
INSERT INTO city (name, state) VALUES ('Семей', 'Восточно-Казахстанская область');
INSERT INTO city (name, state) VALUES ('Талдыкорган', 'Алматинская область');
INSERT INTO city (name, state) VALUES ('Тараз', 'Жамбылская область');
INSERT INTO city (name, state) VALUES ('Темиртау', 'Карагандинская область');
INSERT INTO city (name, state) VALUES ('Туркестан', 'Туркестанская область');
INSERT INTO city (name, state) VALUES ('Уральск', 'Западно-Казахстанская область');
INSERT INTO city (name, state) VALUES ('Усть-Каменогорск', 'Восточно-Казахстанская область');
INSERT INTO city (name, state) VALUES ('Шымкент', 'Южно-Казахстанская область');
INSERT INTO city (name, state) VALUES ('Экибастуз', 'Павлодарская область');
INSERT INTO city (name, state) VALUES ('Жезказган', 'Карагандинская область');
INSERT INTO city (name, state) VALUES ('Жанаозен', 'Мангистауская область');
INSERT INTO city (name, state) VALUES ('Жаркент', 'Алматинская область');
INSERT INTO city (name, state) VALUES ('Жетысай', 'Южно-Казахстанская область');
INSERT INTO city (name, state) VALUES ('Жетысу', 'Алматинская область');
INSERT INTO city (name, state) VALUES ('Жолымбет', 'Актюбинская область');
INSERT INTO city (name, state) VALUES ('Жем', 'Алматинская область');
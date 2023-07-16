-- Создание таблицы campaigns
CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

-- Создание таблицы items
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    campaign_id INT REFERENCES campaigns(id),
    name VARCHAR(255),
    description VARCHAR(255),
    priority INT,
    removed BOOLEAN,
    created_at TIMESTAMP
);

-- Создание индекса на поле campaign_id
CREATE INDEX idx_items_campaign_id ON items(campaign_id);


-- Создание функции, которая будет автоматически устанавливать приоритет
CREATE OR REPLACE FUNCTION set_item_priority()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.priority := (SELECT COALESCE(MAX(priority), 0) + 1 FROM items);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Создание триггера, который будет вызывать функцию set_item_priority() при добавлении записи в таблицу items:
CREATE TRIGGER set_item_priority_trigger
    BEFORE INSERT ON items
    FOR EACH ROW
EXECUTE FUNCTION set_item_priority();


-- Добавление записи по умолчанию в таблицу campaigns
INSERT INTO campaigns (name) VALUES ('Первая запись');
-- Удаление триггера set_item_priority_trigger
DROP TRIGGER IF EXISTS set_item_priority_trigger ON items;

-- Удаление функции set_item_priority
DROP FUNCTION IF EXISTS set_item_priority();

-- Удаление индекса idx_items_campaign_id
DROP INDEX IF EXISTS idx_items_campaign_id;

-- Удаление таблицы items
DROP TABLE IF EXISTS items;

-- Удаление таблицы campaigns
DROP TABLE IF EXISTS campaigns;
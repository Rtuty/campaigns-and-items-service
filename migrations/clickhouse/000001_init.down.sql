-- Удаление материализованного представления item_mv
DROP MATERIALIZED VIEW IF EXISTS items_mv;

-- Удаление таблицы items
DROP TABLE IF EXISTS items;
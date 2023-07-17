-- Создание таблицы item
CREATE TABLE IF NOT EXISTS items (
                                     Id          Int64,
                                     CampaignId  Int64,
                                     Name        String,
                                     Description String,
                                     Priority    Int64,
                                     Removed     UInt8,
                                     EventTime   Timestamp
) ENGINE = MergeTree()
      ORDER BY Id;

-- Создание материализованного представления item_mv для ускорения чтения
CREATE MATERIALIZED VIEW IF NOT EXISTS items_mv
            ENGINE = SummingMergeTree()
                ORDER BY Id
            POPULATE
AS SELECT
       Id,
       CampaignId,
       Name,
       Description,
       Priority,
       Removed,
       EventTime
FROM items;
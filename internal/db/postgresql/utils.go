package db

import (
	"context"

	"cais/internal/entities"

	"github.com/jackc/pgx/v4"
)

// ScanRows сканиурет строки, которые вернул sql запрос, возращая результат в виде пустого интверфейса
func ScanRows(r pgx.Rows, entityType string) ([]interface{}, error) {
	var result []interface{}

	switch entityType {
	case "items":
		for r.Next() {
			var i entities.Item

			if err := r.Scan(&i.Id, &i.CampaignId, &i.Name, &i.Description, &i.Priority, &i.Removed, &i.CreatedAt); err != nil {
				return nil, err
			}

			result = append(result, i)
		}

	case "campaigns":
		for r.Next() {
			var p entities.Campaign

			if err := r.Scan(&p.Id, &p.Name); err != nil {
				return nil, err
			}

			result = append(result, p)
		}
	}

	if err := r.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// convertEntity Конвертирует пустой интерфейс к массиву сущностей items/campaign
func convertEntity[T entities.Item | entities.Campaign](entity []interface{}) ([]T, error) {
	var entities []T

	for _, e := range entity {
		if e1, ok := e.(T); ok {
			entities = append(entities, e1)
		} else {
			return make([]T, 0), convertErr
		}
	}
	return entities, nil
}

// executeItemsQuery выполняет запрос, который возвращает массив сущностей items
func (d *db) executeItemsQuery(ctx context.Context, query string, args ...interface{}) ([]entities.Item, error) {
	t, err := d.client.Begin(ctx)
	if err != nil {
		d.logger.Warningf("execute query method error: %v, transaction status error: %s", err, beginErr)
		return nil, err
	}

	rows, err := t.Query(ctx, query, args...)
	if err != nil {
		d.logger.Warningf("execute query method error: %v, transaction status error: %s", err, execErr)
		return nil, err
	}

	scRows, err := ScanRows(rows, "items")
	if err != nil {
		d.logger.Warningf("execute query method error: %v, transaction status error: %s", err, scanRowErr)
		return nil, err
	}

	items, err := convertEntity[entities.Item](scRows)
	if err != nil {
		d.logger.Warningf("execute query method error: %v, transaction status error: %s", err, convertErr)
		return nil, err
	}

	if err = t.Commit(ctx); err != nil {
		d.logger.Warningf("execute query method error: %v, transaction status error: %s", err, commitErr)
		return nil, err
	}

	return items, nil
}

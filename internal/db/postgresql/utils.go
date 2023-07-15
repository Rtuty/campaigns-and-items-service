package db

import (
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
			return make([]T, 0), convertError
		}
	}
	return entities, nil
}

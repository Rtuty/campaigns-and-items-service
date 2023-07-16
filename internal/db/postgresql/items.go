package db

import (
	"cais/internal/entities"
	"context"
	"github.com/pkg/errors"
)

var (
	beginErr    = errors.New("error when opening a transaction")
	scanRowErr  = errors.New("error when scanning a string after executing an sql query")
	execErr     = errors.New("error when executing sql query")
	commitErr   = errors.New("transaction commit error")
	rollbackErr = errors.New("failed to rollback transaction")

	convertErr   = errors.New("invalid item type in the array")
	duplicateErr = errors.New("object being created already exists in the database")
)

// ItemsHandleCUD обрабатывает методы CREATE, UPDATE, DELETE для сущности item. Данный метод реализует паттерн абстрактная фабрика
func (d *db) ItemsHandleCUD(ctx context.Context, operation string, itm entities.Item) error {
	// Открываем транзакцию
	t, err := d.client.Begin(ctx)
	if err != nil {
		d.logger.Warningf("Items Handle CUD error: %v, transaction status error: %v", err, beginErr)
		return beginErr
	}

	// При редактировании данных в postgresql происходит установка уровня изоляции транзакции
	if _, err = t.Exec(ctx, "set transaction isolation level serializable"); err != nil {
		d.logger.Warningf("Items Handle CUD error: %v, transaction status error: %v", err, execErr)
		return execErr
	}

	// Выполняем CREATE, UPDATE, DELETE для сущности item в зависимости от значения operation
	switch operation {
	case "create":
		var exists bool
		t.QueryRow(ctx, "select exists(select 1 from items where name = ?)", itm.Name).Scan(&exists) // Проверяем существет ли создаваемый объект в базе

		if exists { // Объект найден ->
			if err = t.Rollback(ctx); err != nil { // Закрываем транзакцию
				d.logger.Warningf("items handler create method error: %v, transaction status error: %v", err, rollbackErr)
				return rollbackErr
			}

			// Логируем и возвращем ошибку которая сообщает о предотвращении дубликации сущности
			d.logger.Warningf("items handler create method error: %v", duplicateErr)
			return duplicateErr
		}

		// Если проверка на дубликат пройдена -> добавляем сущность в базу
		if _, err = t.Exec(ctx,
			`insert into items 
			(campaign_id, name, description, priority, removed, created_at)
			values (?, ?, ?, ?, ?, ?)`,
			itm.CampaignId, itm.Name, itm.Description, itm.Priority, itm.Removed, itm.CreatedAt); err != nil {
			d.logger.Warningf("items handler create method error: %v, transaction status error: %v", err, execErr)
			return execErr
		}

	case "update":

	case "delete":

	}

	// Подтверждаем тразакцию, фиксируются изменения
	if err = t.Commit(ctx); err != nil {
		return commitErr
	}

	return nil
}

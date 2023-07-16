package db

import (
	"context"
	"fmt"

	"cais/internal/entities"

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

var getItemsQuery = `select id,
					campaign_id,
					name,
					description,
					priority,
					removed,
					created_at
			 from items`

// GetAllItems возвращает список всех сущностей items из БД
func (d *db) GetAllItems(ctx context.Context) ([]entities.Item, error) {
	return d.executeItemsQuery(ctx, getItemsQuery)
}

// GetItemsByCampaignId возвращает список всех сущностей items из БД, по заданному id компании
func (d *db) GetItemsByCampaignId(ctx context.Context, id string) ([]entities.Item, error) {
	return d.executeItemsQuery(ctx, getItemsQuery+" where campaign_id = ?", id)
}

// ItemsHandleCUD обрабатывает методы CREATE, UPDATE, DELETE для сущности item. Данный метод реализует паттерн абстрактная фабрика
func (d *db) ItemsHandleCUD(ctx context.Context, operation string, itm entities.Item) error {
	// Открываем транзакцию
	t, err := d.client.Begin(ctx)
	if err != nil {
		d.logger.Warningf("Items Handle CUD error: %v, transaction status error: %v", err, beginErr)
		return beginErr
	}

	/* В уровне изоляции SERIALIZABLE все операции чтения и записи блокируются для других транзакций до завершения текущей.
	Это означает, что другие транзакции не могут выполнить чтение или запись в те же данные, с которыми работает текущая транзакция, пока она не будет завершена.
	Данный механизм обеспечивает полную изоляцию от параллельно выполняющихся транзакций. */
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
		if _, err = t.Exec(ctx,
			`update items 
				set campaign_id = ?, name = ?, description = ?, priority = ?, removed = ?, created_at = ?
				where id = ?`,
			itm.CampaignId, itm.Name, itm.Description, itm.Priority, itm.Removed, itm.CreatedAt, itm.Id); err != nil {
			d.logger.Warningf("items handler update method error: %v, transaction status error: %s", err, execErr)
			return execErr
		}

	case "delete":
		if _, err = t.Exec(ctx, "delete from items where id = ?", itm.Id); err != nil {
			d.logger.Warningf("tems handler delete method error: %v, transaction status error: %s", err, execErr)
			return execErr
		}

	default:
		return fmt.Errorf("operation named '%s' not found", operation)
	}

	// Подтверждаем тразакцию, фиксируются изменения
	if err = t.Commit(ctx); err != nil {
		d.logger.Warningf("Items Handle CUD error: %v, transaction status error: %v", err, commitErr)
		return commitErr
	}

	// Устанавливаем дефолтный уровень изоляции, после завершения транзакции
	if _, err = t.Exec(ctx, "set transaction isolation level default"); err != nil {
		d.logger.Warningf("Items Handle CUD error: %v, transaction status error: %v", err, execErr)
		return execErr
	}

	return nil
}

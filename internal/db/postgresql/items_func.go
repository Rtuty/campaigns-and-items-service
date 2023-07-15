package db

import (
	"cais/internal/entities"
	"context"
	"errors"
)

var (
	beginErr   = errors.New("error when opening a transaction")
	scanRowErr = errors.New("error when scanning a string after executing an sql query")
	execErr    = errors.New("error when executing sql query")
	commitErr  = errors.New("transaction commit error")

	convertError = errors.New("invalid item type in the array")
)

// NewItem создает новую сущность item в БД
func (d *db) NewItem(ctx context.Context, itm entities.Item) (id int64, err error) {
	t, err := d.client.Begin(ctx)
	if err != nil {
		return -1, beginErr
	}

	query := `insert into items 
		(campaign_id, item_name, item_description, item_priority, item_removed, item_created_at)
		values (?, ?, ?, ?, ?, ?) returning id`

	var NewItemId int64

	if err = t.QueryRow(ctx, query, &itm.CampaignId, &itm.Name, &itm.Description, &itm.Priority, &itm.Removed, &itm.CreatedAt).Scan(NewItemId); err != nil {
		return -1, scanRowErr
	}

	if err = t.Commit(ctx); err != nil {
		return -1, commitErr
	}

	return NewItemId, nil
}

// UpdateItem обновляет сущность item в БД
func (d *db) UpdateItem(ctx context.Context, itm entities.Item) error {
	t, err := d.client.Begin(ctx)
	if err != nil {
		return beginErr
	}

	query := `update items 
		set campaign_id = ?, item_name = ?, item_description = ?, item_priority = ?, item_removed = ?, item_created_at = ?
		where id = ?`

	if _, err = t.Exec(ctx, query, &itm.CampaignId, &itm.Name, &itm.Description, &itm.Priority, &itm.Removed, &itm.CreatedAt, &itm.Id); err != nil {
		return execErr
	}

	if err = t.Commit(ctx); err != nil {
		return commitErr
	}

	return nil
}

// DeleteItem удаляет сущность item из БД по id
func (d *db) DeleteItem(ctx context.Context, id string) error {
	t, err := d.client.Begin(ctx)
	if err != nil {
		return beginErr
	}

	if _, err = t.Exec(ctx, "delete from items where id = ?", id); err != nil {
		return execErr
	}

	if err = t.Commit(ctx); err != nil {
		return commitErr
	}

	return nil
}

// GetAllItems возвращает список всех сущностей items из БД
func (d *db) GetAllItems(ctx context.Context) ([]entities.Item, error) {
	t, err := d.client.Begin(ctx)
	if err != nil {
		return nil, beginErr
	}
	query := "select id, campaign_id, item_name, item_description, item_priority, item_removed, item_created_at from items"

	rows, err := t.Query(ctx, query)
	if err != nil {
		return nil, execErr
	}

	scRows, err := ScanRows(rows, "items")
	if err != nil {
		return nil, scanRowErr
	}

	items, err := convertEntity[entities.Item](scRows)
	if err != nil {
		return nil, convertError
	}

	if err = t.Commit(ctx); err != nil {
		return nil, commitErr
	}

	return items, nil
}

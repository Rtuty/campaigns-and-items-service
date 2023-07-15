package db

import (
	"cais/internal/entities"
	"cais/pkg/logger"
	"cais/pkg/pgclient"
	"context"
)

type Storage interface {
	NewItem(ctx context.Context, itm entities.Item) (id int64, err error)
	UpdateItem(ctx context.Context, itm entities.Item) error
	DeleteItem(ctx context.Context, id string) error

	GetAllItems(ctx context.Context) ([]entities.Item, error)
}

type db struct {
	client pgclient.Client
	logger *logger.Logger
}

func NewRepository(client pgclient.Client, logger *logger.Logger) Storage {
	return &db{
		client: client,
		logger: logger,
	}
}

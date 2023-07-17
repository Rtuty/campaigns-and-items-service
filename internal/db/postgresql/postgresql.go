package db

import (
	"cais/internal/entities"
	"cais/pkg/logger"
	"cais/pkg/pgclient"
	"context"
)

type Storage interface {
	ItemsHandleCUD(ctx context.Context, operation string, itm entities.Item) error

	GetAllItems(ctx context.Context) ([]entities.Item, error)
	GetItemsByCampaignId(ctx context.Context, id int64) ([]entities.Item, error)
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

package db

import (
	"cais/internal/entities"
	"cais/pkg/clickapi"
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
	clickh clickapi.ClickHouse
}

func NewRepository(client pgclient.Client, logger *logger.Logger, clickh clickapi.ClickHouse) Storage {
	return &db{
		client: client,
		logger: logger,
		clickh: clickh,
	}
}

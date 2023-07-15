package db

import (
	"cais/pkg/logger"
	"cais/pkg/pgclient"
)

type Storage interface{}

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

package app

import (
	"context"
	"log"

	"cais/internal/api"
	"cais/internal/db/postgresql"
	"cais/pkg/logger"
	"cais/pkg/pgclient"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

// Иннициализация переменных окружения
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func RunServiceInstance() {
	ctx := context.Background()
	log := logger.GetLogger()

	pgDataSource, err := pgclient.GetDataConnection()
	if err != nil {
		log.Printf("get data connection error: %v", err)
	}

	pgcli, err := pgclient.NewClient(ctx, 5, pgDataSource, log)
	if err != nil {
		log.Printf("get new postgresql client error: %v", err)
	}

	var repo = db.NewRepository(pgcli, log)

	var h = api.NewItemHandler(ctx, repo, log)
	r := httprouter.New()

	r.GET("/items/all", h.GetAllItems)
}

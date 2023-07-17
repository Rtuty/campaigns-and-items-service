package app

import (
	"cais/pkg/clickapi"
	"context"
	"log"
	"net/http"

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
	// Создаем контекст и логгер
	ctx := context.Background()
	log := logger.GetLogger()

	// Получаем данные для подключения базе данных
	pgDataSource, err := pgclient.GetDataConnection()
	if err != nil {
		log.Printf("get data connection error: %v", err)
	}

	// Подключаемся к postgresql
	log.Info("getting postgresql client")
	pgcli, err := pgclient.NewClient(ctx, 5, pgDataSource, log)
	if err != nil {
		log.Printf("get new postgresql client error: %v", err)
	}

	defer pgcli.Close()

	// Миграция up postgresql (встраивается в бинарник исполняемого файла)
	if err := pgclient.PostgresMigration(pgDataSource, "up"); err != nil {
		log.Printf("get new postgresql client error: %v", err)
	}

	// Получаем подключение к clickhouse
	clHs, err := clickapi.ClickHouseConnection(ctx, log)
	if err != nil {
		log.Printf("get new clickhouse connection error: %v", err)
	}

	// Получаем интерфейс repository, который реализует функционал для работы с сущностями в среде postgresql
	var repo = db.NewRepository(pgcli, log, clHs)

	// Обработчик для роутера
	var h = api.NewItemHandler(ctx, repo, log)
	r := httprouter.New()

	r.GET("/items/get/all", h.GetAllItems)
	r.GET("/items/get/:id", h.GetItemsByCampaignId)

	r.POST("/items/new", h.CreateNewItem)
	r.PATCH("items/update", h.UpdateItem)
	r.DELETE("items/delete/:id", h.DeleteItem)

	// Запуск сервера
	log.Info("starting server on 8080...")
	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Infof("failed to start the server, error: %v", err)
	}
}

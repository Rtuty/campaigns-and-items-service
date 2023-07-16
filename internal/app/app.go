package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"

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

	// Выполняем команду make из Makefile для запуска миграции postgres
	var mgrCmd = exec.Command("make", "migrate-up")
	mgrCmd.Stderr = os.Stderr
	mgrCmd.Stdout = os.Stdout

	if err := mgrCmd.Run(); err != nil {
		log.Infof("migration command execution error: %v", err)
	}

	// Получаем интерфейс repository, который реализует функционал для работы с сущностями в среде postgresql
	var repo = db.NewRepository(pgcli, log)

	// Обработчик для роутера
	var h = api.NewItemHandler(ctx, repo, log)
	r := httprouter.New()

	r.GET("/items/all", h.GetAllItems)
	r.POST("/items/new", h.CreateNewItem)
	r.PATCH("items/update", h.UpdateItem)
	r.DELETE("items/delete/:id", h.DeleteItem)

	// Запуск сервера
	log.Info("starting server on 8080...")
	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Infof("failed to start the server, error: %v", err)
	}
}

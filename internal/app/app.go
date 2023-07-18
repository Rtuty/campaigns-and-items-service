package app

import (
	"cais/pkg/caching"
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"

	"cais/internal/api"
	"cais/internal/db/postgresql"
	"cais/pkg/clickapi"
	"cais/pkg/logger"
	"cais/pkg/msgbroker"
	"cais/pkg/pgclient"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/nats.go"
)

// Иннициализация переменных окружения
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// InstanceOpts - Настройки для запуска инстанса (заполняется в main)
type InstanceOpts struct {
	Port      string
	CfgPgVars []string
	RedisVar  caching.RedisEnvars
}

func RunServiceInstance(opts InstanceOpts) {
	pgCh := make(chan *pgxpool.Pool)
	rdCh := make(chan *redis.Client)
	clCh := make(chan driver.Conn)

	// Создаем контекст и логгер
	ctx := context.Background()
	log := logger.GetLogger()

	go func() {
		// Получаем данные для подключения базе данных
		log.Info("get data for postgresql connection")
		pgDataSource, err := pgclient.GetPgDataConnection(opts.CfgPgVars)
		if err != nil {
			log.Warningf("get data connection error: %v", err)
		}

		// Получение клиента postgresql
		log.Info("getting postgresql client")
		pgcli, err := pgclient.NewPostgresClient(ctx, 5, pgDataSource, log)
		if err != nil {
			log.Warningf("get new postgresql client error: %v", err)
		}

		// Миграция up postgresql (встраивается в бинарник исполняемого файла)
		log.Info("start postgresql migration up")
		if err := pgclient.PostgresMigration(pgDataSource, "up"); err != nil {
			log.Warningf("get new postgresql client error: %v", err)
		}

		pgCh <- pgcli
	}()

	go func() {
		// Получаем клиент redis
		log.Info("getting redis client")
		rdCl, err := caching.NewRedisClient(ctx, log, opts.RedisVar)
		if err != nil {
			log.Warningf("redis client connection error: %v", err)
		}

		rdCh <- rdCl
	}()

	go func() {
		// Получаем подключение к clickhouse
		log.Info("getting click house client")
		clHs, err := clickapi.ClickHouseConnection(ctx, log)
		if err != nil {
			log.Warningf("get new clickhouse connection error: %v", err)
		}

		clCh <- clHs
	}()

	// Ожидание результатов подключений
	pgcli := <-pgCh
	rdCl := <-rdCh
	clHs := <-clCh

	defer func() { // Закрытие postgres, redis, clickhouse клиентов, после завершения программы
		pgcli.Close()
		clHs.Close()
		rdCl.Close()
	}()

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

	// NATS run and serve
	msgbroker.ConnectAndServeNATS(ctx, log, pgcli, nats.DefaultURL, "items")

	// Запуск сервера
	log.Infof("starting server on %s...", opts.Port)
	if err := http.ListenAndServe(opts.Port, r); err != nil {
		log.Infof("failed to start the server, error: %v", err)
	}
}

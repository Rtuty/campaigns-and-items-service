package msgbroker

import (
	"cais/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/nats.go"
)

func ConnectAndServeNATS(ctx context.Context, log *logger.Logger, pg *pgxpool.Pool, conUrl string, listenChannel string) {
	nc, err := nats.Connect(conUrl) // Создание соединения с сервером Nats
	if err != nil {
		log.Fatalf("msgbroker doesn't connected: %v", err)
	}

	defer nc.Close()

	conn, err := pg.Acquire(ctx) // Подписка на события изменения записей в базе данных
	if err != nil {
		log.Warning(err)
	}

	defer conn.Release()

	_, err = conn.Exec(ctx, fmt.Sprintf("LISTEN %s;", listenChannel))
	if err != nil {
		log.Fatal(err)
	}

	// Обработка событий изменения записей
	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			log.Fatal(err)
		}

		// Отправка лога в Clickhouse через очередь Nats
		if err := nc.Publish("items", []byte(notification.Payload)); err != nil {
			log.Fatal(err)
		}
	}
}

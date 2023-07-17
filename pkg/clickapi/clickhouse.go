package clickapi

import (
	"cais/pkg/logger"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// ClickHouse интерфейс реализует функционал работы с click house
type ClickHouse interface {
	Contributors() []string
	ServerVersion() (*driver.ServerVersion, error)
	Select(ctx context.Context, dest any, query string, args ...any) error
	Query(ctx context.Context, query string, args ...any) (driver.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) driver.Row
	PrepareBatch(ctx context.Context, query string) (driver.Batch, error)
	Exec(ctx context.Context, query string, args ...any) error
	AsyncInsert(ctx context.Context, query string, wait bool) error
	Ping(context.Context) error
	Stats() driver.Stats
	Close() error
}

// ClickHouseConnection возвращает подключение к базе данных clickhouse
func ClickHouseConnection(ctx context.Context, log *logger.Logger) (driver.Conn, error) {
	var (
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{"<CLICKHOUSE_SECURE_NATIVE_HOSTNAME>:9440"},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: "default",
				Password: "<DEFAULT_USER_PASSWORD>",
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},

			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	)
	if err != nil {
		log.Error("open clickhouse error: %v", err)
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil { // проверка работает ли clickhouse
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Infof("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		log.Errorf("ping clickhouse error: %v", err)
		return nil, err
	}

	log.Info("connection with clickhouse is established")
	return conn, nil
}

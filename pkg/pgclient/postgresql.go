package pgclient

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"cais/pkg/logger"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// Client интерфейс реализует функционал работы с postgres
type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) // Исполнение запроса (update, delete, etc)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)              // Выполнить запрос, который вернет множество строк (>1)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row                     // Вернуть одну строку из запроса
	Begin(ctx context.Context) (pgx.Tx, error)                                                 // Открытие и получение сущности транзакции postgres
}

/*
FAQ: Как работает подключение к базе данных postgresql

	GetDataConnection:
		1. Достаем переменные окружения из файла .env
		2. Если все нужные переменные существуют и иметют значение -> записываем в структуру и возвращаем данные для подключения

	NewClient:
		1. Отправляем в аргументы контекст программы, максимальное количество попыток для подключения, данные из .env и логгер
		2. Если подключение к postgresql произведенено корректно, то возвращаем клиент для работы с базой данных
*/

type dataSource struct{ Host, Port, User, Passwd, Dbname string }

func GetDataConnection() (dataSource, error) {
	var postgresCon = dataSource{}

	envVars := []string{"HOST", "PORT", "USER", "PASSWD", "DBNAME"} // Данные, которые хотим получить из переменных окружения

	for _, v := range envVars {
		value := os.Getenv(v)
		if value == "" {
			return dataSource{}, errors.Errorf("invalid environment variable %s", v)
		} else { // Если переменная окружения найдена - записываем результат в postgresCon
			field := reflect.ValueOf(&postgresCon).Elem().FieldByNameFunc(
				func(fieldName string) bool {
					return strings.EqualFold(fieldName, v)
				})
			if field.IsValid() {
				field.SetString(value)
			}
		}
	}

	return postgresCon, nil
}

func NewClient(ctx context.Context, maxAttempts int, ds dataSource, log *logger.Logger) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", ds.User, ds.Passwd, ds.Host, ds.Port, ds.Dbname)
	err = doWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}

package pgclient

import (
	"cais/migrations"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"
)

func PostgresMigration(ds dataSource, cmd string) error {
	d, err := iofs.New(migrations.FS, "postgres")
	if err != nil {
		log.Fatalf("migration command execution error: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		ds.User, ds.Passwd, ds.Host, ds.Port, ds.Dbname,
	))

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange { // если нет изменений, то пропускаем миграцию (не ошибка)
			return fmt.Errorf("migration up error: %v", err)
		}

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange { // если нет изменений, то пропускаем миграцию (не ошибка)
			return fmt.Errorf("migration up error: %v", err)
		}

	default:
		return fmt.Errorf("migration command '%s' is incorrect", cmd)
	}

	return nil
}

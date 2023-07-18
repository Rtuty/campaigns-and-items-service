include .env # Выгружаем данные для миграции из .env
export

run:
	go mod download
	go run ./cmd/campaigns-and-items-service/main.go

migrate-up:
	migrate -path ./migrate/postgres -database 'postgres://postgres:$(PG_PASSWD)@$(PG_HOST):$(PG_PORT)/$(PG_DBNAME)?sslmode=$(PG_SSLMODE)' up

migrate-down:
	migrate -path ./migrate/postgres -database 'postgres://postgres:$(PG_PASSWD)@$(PG_HOST):$(PG_PORT)/$(PG_DBNAME)?sslmode=$(PG_SSLMODE)' down

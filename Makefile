include .env # Выгружаем данные для миграции из .env
export

run:
	go mod download
	go run ./cmd/campaigns-and-items-service/main.go

migrate-up:
	migrate -path ./migrate/postgres -database 'postgres://postgres:$(PASSWD)@$(HOST):$(PORT)/$(DBNAME)?sslmode=$(SSLMODE)' up

migrate-down:
	migrate -path ./migrate/postgres -database 'postgres://postgres:$(PASSWD)@$(HOST):$(PORT)/$(DBNAME)?sslmode=$(SSLMODE)' down

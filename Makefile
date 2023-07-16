include .env # Выгружаем данные для миграции из .env
export

migrate-up:
	migrate -path ./migrate -database 'postgres://postgres:$(PASSWD)@$(HOST):$(PORT)/$(DBNAME)?sslmode=$(SSLMODE)' up

migrate-down:
	migrate -path ./migrate -database 'postgres://postgres:$(PASSWD)@$(HOST):$(PORT)/$(DBNAME)?sslmode=$(SSLMODE)' down
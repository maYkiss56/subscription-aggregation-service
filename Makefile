include .env
export

APP=bin/sas
MIGRATIONS_DIR=migrations

.PHONY: build b
.PHONY: run r
.PHONY: test t
.PHONY: clean c
.PHONY: migrate-new mn
.PHONY: migrate-up mu
.PHONY: migrate-down md
.PHONY: migrate-drop mdr
.PHONY: migrate-status ms
.PHONY: swag-init si
.PHONY: swag-clean sc
.PHONY: help h

build:
	go build -o $(APP) ./cmd/main.go

b: build

run: build
	./$(APP)

r: run

test:
	go test -v ./...
t: test

clean:
	rm -rf ./bin || true

c: clean

migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

mn: migrate-new

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) up

mu: migrate-up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) down

md: migrate-down

migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) drop

mdr: migrate-drop

migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database $(DB_URL) version

ms: migrate-status

swag-init:
	swag init --parseDependency -g ./cmd/main.go -o ./docs

si: swag-init

swag-clean:
	rm -rf docs/

sc: swag-clean

help:
	@echo "Available commands:"
	@echo " make build           - Build the application"
	@echo " make run             - Build and run the application"
	@echo " make test            - Run tests"
	@echo " make clean           - Remove the compiled binary"
	@echo " make migrate-new     - Create a new migration"
	@echo " make migrate-up      - Apply all up migrations"
	@echo " make migrate-down    - Roll back the last migration"
	@echo " make migrate-drop    - Drop all migrations"
	@echo " make migrate-status  - Show current migration version"
	@echo " make swag-init       - Generate Swagger documentation"
	@echo " make swag-clean      - Remove generated Swagger files"
	@echo " make help            - Show this help message"

h: help

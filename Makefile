include .env

# THIS FILE IS ONLY SUITABLE FOR LOCAL DEVELOPMENT - WIP

build:
	go build -o main ./cmd/bot/main.go

run:
	go run ./cmd/bot/main.go

createMigration:
	migrate create -ext sql -dir ./data/migrations/ -seq $(name)

createDump:
	DB_USER=${DB_USER} DB_NAME=${DB_NAME} sh ./scripts/db/create_dump.sh

applyDump:
	DB_USER=${DB_USER} DB_NAME=${DB_NAME} sh ./scripts/db/find_apply_dump.sh

migrationDown: createDump
	migrate -path ./data/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose \
	down

_migrationUp:
	migrate -path ./data/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose \
	up

migrationUp: _migrationUp applyDump

.PHONY: build run migrationUp migrationDown applyDump createMigration

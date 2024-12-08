include .env
include ./scripts/lts_dump.sh

build:
	go build -o main ./cmd/bot/main.go

run:
	go run ./cmd/bot/main.go

createMigration:
	migrate create -ext sql -dir ./data/migrations/ -seq init_schema

createDump:
	docker exec -it airquality-db-container sh -c "pg_dump -U ${DB_USER} ${DB_NAME} > dump_$(date +%Y-%m-%d_%H-%M-%S).sql;"

applyDump:
	find_lts_dump

migrationDown: createDump
	migrate -path ./data/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose \
	down

_migrationUp:
	migrate -path ./data/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose \
	up

migrationUp: _migrationUp applyDump

.PHONY: build run migrationUp migrationDown applyDump createMigration

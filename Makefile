include .env

# THIS FILE IS ONLY SUITABLE FOR LOCAL DEVELOPMENT - WIP
build:
	go build -o main ./cmd/bot/main.go

run:
	go run ./cmd/bot/main.go

create-migration:
	@bash -c 'read -p "Please provide migration name: " name && \
	echo $$name && \
	migrate create -ext sql -dir ./data/migrations/ -seq $$name'

create-dump:
	DB_USER=${DB_USER} DB_NAME=${DB_NAME} sh ./scripts/db/create_dump.sh

apply-dump:
	DB_USER=${DB_USER} DB_NAME=${DB_NAME} sh ./scripts/db/find_apply_dump.sh

generate-coverage:
	go test -v -coverprofile ./tmp/coverage.out ./internal/... && \
	go tool cover -html=./tmp/coverage.out -o ./tmp/coverage.html && \
	google-chrome ./tmp/coverage.html

migration-down: create-dump
	migrate -path ./data/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose \
	down

.PHONY: build run migration-down apply-dump create-migration generate-html-coverage

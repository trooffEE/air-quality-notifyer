include .env

build:
	go build -o main ./cmd/bot/main.go

run:
	go run ./cmd/bot/main.go

applyDump:
	docker cp ./data/dump/dump.sql airquality-db-container:/dump.sql && \
	docker exec -it airquality-db-container sh -c "psql -U ${DB_USER} ${DB_NAME} < dump.sql;"

migrationDown:
	migrate -path ./data/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose down

migrationUp:
	migrate -path ./data/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose up

.PHONY: build run applyDump migrationDown migrationUp

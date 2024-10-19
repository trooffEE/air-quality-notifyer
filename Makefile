.PHONY: run
run:
	go run ./cmd/bot/main.go

.PHONY: build
build:
	go build -o main ./cmd/bot/main.go

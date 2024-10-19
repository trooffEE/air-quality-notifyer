.PHONY: run
run:
	go run ./main.go

.PHONY: build
build:
	go build -o main ./main.go

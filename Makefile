.PHONY: tidy build run sqlc migrate-up migrate-down

POSTGRES_URL ?= postgres://postgres:postgres@localhost:5432/ai_chat?sslmode=disable

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

sqlc:
	sqlc generate

migrate-up:
	migrate -path migrations -database "$(POSTGRES_URL)" up

migrate-down:
	migrate -path migrations -database "$(POSTGRES_URL)" down 1

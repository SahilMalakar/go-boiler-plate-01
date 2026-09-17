.PHONY: build run clean migrate-up migrate-down

build:
	@go build -o bin/api ./cmd/api

run: build
	@./bin/api

clean:
	@rm -rf bin

migrate-up:
	@go run ./cmd/migrate up
	
migrate-down:
	@go run ./cmd/migrate down


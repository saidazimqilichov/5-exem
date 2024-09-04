gen-hotel:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	./protos/hotel/hotel.proto




run:
	@go run cmd/main.go

CURRENT_DIR := $(shell pwd)
DB_URL := "postgres://postgres:7777@localhost:5432/hotel?sslmode=disable"

migrate-up:
	migrate -path migrations/ -database $(DB_URL) up


migrate-down:
	migrate -path migrations/ -database  $(DB_URL) down

migrate-force:
	migrate -path migrations/ -database  $(DB_URL) force 1

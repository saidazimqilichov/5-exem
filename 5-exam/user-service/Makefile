tidy:
	@go mod tidy 

run:
	@go run cmd/main.go 

migrate:
	migrate create -dir ./internal/db -ext sql db

migrate-up:
	migrate -path ./internal/db -database "postgres://postgres:7777@localhost:5432/real?sslmode=disable" up

migrate-down:
	migrate -path ./internal/db -database "postgres://postgres:7777@localhost:5432/real?sslmode=disable" down

migrate-force:
	migrate -path ./internal/db -database "postgres://postgres:7777@localhost:5432/real?sslmode=disable" force 

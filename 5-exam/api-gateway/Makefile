gen-hotel:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	./proto/hotel.proto



gen-user:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	./proto/user.proto


gen-booking:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	./protos/hotel/hotel.proto

# run:
# 	@go run ./cmd/main.go


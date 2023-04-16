# Генерация сервера
grpc_server:
	protoc --go_out=./internal/api/grpc/gen --proto_path=./internal/api/grpc/proto --go_opt=paths=source_relative --go-grpc_out=./internal/api/grpc/gen --go-grpc_opt=require_unimplemented_servers=false --go-grpc_opt=paths=source_relative $(shell find ./internal/api/grpc/proto -iname "*.proto")

# Запуск приложения
run:
	go run ./cmd/main.go --config=./config.yml
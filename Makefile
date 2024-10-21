include .env

LOCAL_BIN:=$(CURDIR)/bin
MIGRATION_DIR=./database/migrations
MIGRATION_DSN="host=$(DB_HOST) port=$(DB_PORT) dbname=$(DB_DATABASE) user=$(DB_USER) password=$(DB_PASSWORD)"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.22.1

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-user-api

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path grpc/api/user/v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	grpc/api/user/v1/user.proto

migration-status:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) status -v

migration-migrate:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) up -v

migration-rollbacK:
	$(LOCAL_BIN)/goose -dir $(MIGRATION_DIR) postgres $(MIGRATION_DSN) down -v

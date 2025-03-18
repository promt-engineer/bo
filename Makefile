goose-install:
	go get -u github.com/pressly/goose/v3/cmd/goose

migration:
	goose -dir ./migrations create $(name) sql

migrate-up:
	goose -dir ./migrations -table goose_db_version postgres $(LOCAL_DB_URL) up

migrate-down:
	goose -dir ./migrations -table goose_db_version postgres $(LOCAL_DB_URL) down

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
 			--go-grpc_opt=require_unimplemented_servers=false \
 			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
 			./pkg/backoffice/main.proto

proto-history:
	protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_opt=require_unimplemented_servers=false \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				./pkg/history/main.proto

swag:
	swag init -g cmd/backoffice/main.go

lint:
	golangci-lint run -v
MIGRATION_NAME=init_schema
DB_URL=postgresql://root:secret@localhost:5432/go_exchange?sslmode=disable

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## migrate_create: create migration files with name MIGRATION_NAME
migrate_create:
	migrate create -ext sql -dir db/migration -seq ${MIGRATION_NAME}

## migrate_up: migrate db for the last version
migrate_up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

## migrate_down: migrate db for the first version
migrate_down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

## migrate_drop: drop averything inside database
migrate_drop:
	migrate -path db/migration -database "$(DB_URL)" -verbose drop -f

## sqlc: generate go code from SQL
sqlc:
	sqlc generate

## test: execute tests
test:
	go test -v -cover ./...

## server: starts server
server:
	go run main.go

## mock: generates mock interfaces in reflect mode
mock:
	mockgen -package mockdb -destination db/mock/store.go go-exchange/db/sqlc Store

## proto: generates go code with protobuf
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=go_exchange \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

## db_docs: generates database documentation
# To set password: $ dbdocs password --set my_password --project my_project_name
db_docs:
	dbdocs build doc/db.dbml

## db_schema: generates database schema from dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: up up_build down \
		migrate_create migrate_up migrate_down migrate_drop \
		sqlc test server mock proto \
		db_docs db_schema

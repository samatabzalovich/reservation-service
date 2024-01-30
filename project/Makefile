
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
PRODUCT_BINARY=productServiceApp
MAIL_BINARY=mailerApp
LISTENER_BINARY=listenerApp


## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"
#build_logger build_mail build_listener build_product
## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth
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

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

build_product:
	@echo "Building product binary..."
	cd ../product-service && env GOOS=linux CGO_ENABLED=0 go build -o ${PRODUCT_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"

AUTH_PROTO_DIR := ../authentication-service/auth_proto
AUTH_PROTO_FILE := $(AUTH_PROTO_DIR)/auth.proto
AUTH_OUT_DIR := $(AUTH_PROTO_DIR)

BROKER_AUTH_PROTO_DIR := ../broker-service/proto_files/auth
BROKER_AUTH_PROTO_FILE := $(BROKER_AUTH_PROTO_DIR)/auth.proto
BROKER_AUTH_OUT_DIR := $(BROKER_AUTH_PROTO_DIR)

BROKER_INST_PROTO_DIR := ../broker-service/proto_files/institution_proto
BROKER_INST_PROTO_FILE := $(BROKER_INST_PROTO_DIR)/inst.proto
BROKER_INST_OUT_DIR := $(BROKER_INST_PROTO_DIR)

INSTITUTION_PROTO_DIR := ../institution-service/proto_files/institution_proto
INSTITUTION_PROTO_FILE := $(INSTITUTION_PROTO_DIR)/inst.proto
INSTITUTION_OUT_DIR := $(INSTITUTION_PROTO_DIR)


## build_proto: builds the proto_files files

build_proto:
	@echo "Building proto files..."
	protoc --go_out=${AUTH_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${AUTH_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${AUTH_PROTO_DIR} ${AUTH_PROTO_FILE}
	protoc --go_out=${BROKER_AUTH_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${BROKER_AUTH_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${BROKER_AUTH_PROTO_DIR} ${BROKER_AUTH_PROTO_FILE}
	protoc --go_out=${INSTITUTION_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${INSTITUTION_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${INSTITUTION_PROTO_DIR} ${INSTITUTION_PROTO_FILE}
	protoc --go_out=${BROKER_INST_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${BROKER_INST_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${BROKER_INST_PROTO_DIR} ${BROKER_INST_PROTO_FILE}
	@echo "Done!"


## Migration paths
DSN := postgres://postgres:190704Samat@localhost/postgres?sslmode=disable
INSTITUTION_DIR_MIGR := ../institution-service/migrations
AUTH_DIR_MIGR := ../authentication-service/migrations

migrate_up: 
	@echo "Migrating up..."
	migrate -path ${INSTITUTION_DIR_MIGR} -database "${DSN}" up
	migrate -path ${AUTH_DIR_MIGR} -database "${DSN}" up
	@echo "Done!"

migrate_down:
	@echo "Migrating down..."
	migrate -path ${INSTITUTION_DIR_MIGR} -database "${DSN}" down
	migrate -path ${AUTH_DIR_MIGR} -database "${DSN}" down
	@echo "Done!"
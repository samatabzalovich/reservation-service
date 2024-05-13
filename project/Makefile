
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
INSTITUTION_BINARY=instApp
STAFF_BINARY=staffApp
APPOINTMENT_BINARY=appointmentApp
NOTIFICATION_BINARY=notificationApp
QUEUE_BINARY=queueApp


build_only: build_broker build_auth build_inst build_staff build_appointment build_notify build_queue

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"
#build_logger build_mail build_listener build_product
## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_inst build_staff build_appointment build_notify build_queue
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

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_inst : builds the institution binary as a linux executable
build_inst:
	@echo "Building institution binary..."
	cd ../institution-service && env GOOS=linux CGO_ENABLED=0 go build -o ${INSTITUTION_BINARY} ./cmd/api
	@echo "Done!"

build_staff:
	@echo "Building staff binary..."
	cd ../staff-service && env GOOS=linux CGO_ENABLED=0 go build -o ${STAFF_BINARY} ./cmd/api
	@echo "Done!"

build_appointment:
	@echo "Building appointment binary..."
	cd ../appointment-service && env GOOS=linux CGO_ENABLED=0 go build -o ${APPOINTMENT_BINARY} ./cmd/api
	@echo "Done!"

build_notify:
	@echo "Building notification binary..."
	cd ../notification-service && env GOOS=linux CGO_ENABLED=0 go build -o ${NOTIFICATION_BINARY} ./cmd/api
	@echo "Done!"

build_queue:
	@echo "Building queue binary..."
	cd ../queue-management-service && env GOOS=linux CGO_ENABLED=0 go build -o ${QUEUE_BINARY} ./cmd/api
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

AUTH_EMP_PROTO_DIR := ../authentication-service/employee_proto
AUTH_EMP_PROTO_FILE := $(AUTH_EMP_PROTO_DIR)/employee_proto.proto
AUTH_EMP_OUT_DIR := $(AUTH_EMP_PROTO_DIR)

EMPLOYEE_PROTO_DIR := ../staff-service/proto_files/employee_proto
EMPLOYEE_PROTO_FILE := $(EMPLOYEE_PROTO_DIR)/employee_proto.proto
EMPLOYEE_OUT_DIR := $(EMPLOYEE_PROTO_DIR)

EMPLOYEE_AUTH_PROTO_DIR := ../staff-service/proto_files/auth_proto
EMPLOYEE_AUTH_PROTO_FILE := $(EMPLOYEE_AUTH_PROTO_DIR)/auth.proto
EMPLOYEE_AUTH_OUT_DIR := $(EMPLOYEE_AUTH_PROTO_DIR)

EMPLOYEE_INST_PROTO_DIR := ../staff-service/proto_files/institution_proto
EMPLOYEE_INST_PROTO_FILE := $(EMPLOYEE_INST_PROTO_DIR)/inst.proto
EMPLOYEE_INST_OUT_DIR := $(EMPLOYEE_INST_PROTO_DIR)

INSTITUTION_PROTO_DIR := ../institution-service/proto_files/institution_proto
INSTITUTION_PROTO_FILE := $(INSTITUTION_PROTO_DIR)/inst.proto
INSTITUTION_OUT_DIR := $(INSTITUTION_PROTO_DIR)

APPOINTMENT_AUTH_PROTO_DIR := ../appointment-service/proto_files/auth_proto
APPOINTMENT_AUTH_PROTO_FILE := $(APPOINTMENT_AUTH_PROTO_DIR)/auth.proto
APPOINTMENT_AUTH_OUT_DIR := $(APPOINTMENT_AUTH_PROTO_DIR)



## build_proto: builds the proto_files files
# export GO_PATH=~/go
# export PATH=$PATH:/$GO_PATH/bin
build_proto:
	@echo "Building proto files..."
	protoc --go_out=${AUTH_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${AUTH_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${AUTH_PROTO_DIR} ${AUTH_PROTO_FILE}
	protoc --go_out=${BROKER_AUTH_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${BROKER_AUTH_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${BROKER_AUTH_PROTO_DIR} ${BROKER_AUTH_PROTO_FILE}
	protoc --go_out=${INSTITUTION_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${INSTITUTION_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${INSTITUTION_PROTO_DIR} ${INSTITUTION_PROTO_FILE}
	protoc --go_out=${BROKER_INST_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${BROKER_INST_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${BROKER_INST_PROTO_DIR} ${BROKER_INST_PROTO_FILE}
	protoc --go_out=${AUTH_EMP_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${AUTH_EMP_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${AUTH_EMP_PROTO_DIR} ${AUTH_EMP_PROTO_FILE}
	protoc --go_out=${EMPLOYEE_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${EMPLOYEE_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${EMPLOYEE_PROTO_DIR} ${EMPLOYEE_PROTO_FILE}
	protoc --go_out=${EMPLOYEE_AUTH_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${EMPLOYEE_AUTH_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${EMPLOYEE_AUTH_PROTO_DIR} ${EMPLOYEE_AUTH_PROTO_FILE}
	protoc --go_out=${EMPLOYEE_INST_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${EMPLOYEE_INST_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${EMPLOYEE_INST_PROTO_DIR} ${EMPLOYEE_INST_PROTO_FILE}
	protoc --go_out=${APPOINTMENT_AUTH_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${APPOINTMENT_AUTH_OUT_DIR} --go-grpc_opt=paths=source_relative --proto_path=${APPOINTMENT_AUTH_PROTO_DIR} ${APPOINTMENT_AUTH_PROTO_FILE}
	@echo "Done!"


## Migration paths
DSN := postgres://reservehub:190704Samat@reserve-hub-database.cnwmuwismjtn.us-east-1.rds.amazonaws.com:5432/postgres
DIR_MIGRATIONS := ../migrations
NAME := create_queue_table

migrate_up:
	@echo "Migrating up..."
	migrate -path ${DIR_MIGRATIONS} -database "${DSN}" up
	@echo "Done!"



migrate_down:
	@echo "Migrating down..."
	migrate -path ${DIR_MIGRATIONS} -database "${DSN}" down
	@echo "Done!"



migrate_create:
	@echo "Creating migration..."
	migrate create -ext sql -dir ${DIR_MIGRATIONS} -seq ${NAME}
	@echo "Done!"


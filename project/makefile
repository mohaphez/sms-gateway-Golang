BROKER_BINARY=brokerApp
AUTHENTICATION_BINERY=authApp
SMS_BINERY=smsApp
LOGGER_BINARY=loggerServiceApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose  --env-file .env up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_sms build_logger
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose --env-file .env up --build -d
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
	@echo "Building authentication binary..."
	cd ../authentication-service/cmd/api && env GOOS=linux CGO_ENABLED=0 go build -o ../../${AUTHENTICATION_BINERY}
	@echo "Done!"

## build_sms: builds the sms binary as a linux executable
build_sms:
	@echo "Building sms binary..."
	cd ../sms-service/cmd/api && env GOOS=linux CGO_ENABLED=0 go build -o ../../${SMS_BINERY}
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"
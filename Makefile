.PHONY: generate lint-proto build run test test-cover lint docker-build docker-up docker-down

generate: ### generate Go code from proto files
	buf generate

lint-proto: ### lint proto files
	buf lint

build: ### build the application binary
	go build -o ./bin/app ./cmd/app

run: ### run the application locally
	go run ./cmd/app

test: ### run unit tests
	go test ./...

test-cover: ### run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint: ### run golangci-lint
	golangci-lint run ./...

docker-build: ### build docker image
	docker build -f deploy/Dockerfile -t card-validator .

docker-up: ### start with docker compose
	docker compose -f deploy/docker-compose.yml up -d

docker-down: ### stop docker compose
	docker compose -f deploy/docker-compose.yml down

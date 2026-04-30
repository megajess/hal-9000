## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run: run the server
.PHONY: run
run:
	go run ./...

## test/verbose: run unit tests with verbose output
.PHONY: test/verbose
test/verbose:
	go test -v ./...

## build: build the server binary
.PHONY: build
build:
	go build -o hal ./...

## test: run unit tests
.PHONY: test
test:
	go test ./...

## db/migrations/new name=$1: create a new migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database "sqlite3://${HAL_DB_PATH}" up

## db/migrations/down: apply all down migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations -database "sqlite3://${HAL_DB_PATH}" down
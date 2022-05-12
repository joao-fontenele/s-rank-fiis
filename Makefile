TAG=$(shell git rev-parse --short HEAD)
GOCMD?=$(shell which go)
#GOCMD=go1.18

init: install-deps
	touch .env

.PHONY: install-deps
install-deps:
	${GOCMD} install github.com/cespare/reflex@v0.3

.PHONY: run
run: tidy start-db
	${GOCMD} run ./cmd/server/main.go

.PHONY: start-db
start-db:
	docker-compose up -d

.PHONY: run-watch
run-watch:
	reflex -d none -s -r '\.go$$' make run

.PHONY: build
build:
	go build -ldflags "-X main.Version=${TAG}" -o ./build/server ./cmd/server/main.go

.PHONY: test
test:
	${GOCMD} test ./...

.PHONY: test-watch
test-watch:
	reflex -d none -s -r '\.go$$' make test

.PHONY: tidy
tidy:
	${GOCMD} mod tidy

.PHONY: psql
psql:
	docker-compose exec db /bin/sh -c "PGPASSWORD=pwd psql -U usr ranks"

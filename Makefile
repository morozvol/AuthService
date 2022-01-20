.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: migrate
migrate:
	migrate -path ./migrations -database 'postgres://postgres:$(password)@localhost:5430/auth?sslmode=disable' up

create_migrate:
	migrate create -ext sql -dir ./migrations $(migrate_name)

.DEFAULT_GOAL := build

.PHONY:
.SILENT:

build:
	go build -o ./bin/main  github.com/morozvol/AuthService/cmd/authserver
test:
	go test -v -race -timeout 30s github.com/morozvol/AuthService/cmd/authserver

migrate:
	migrate -path ./db/migrations -database 'postgres://postgres:$(password)@localhost:5430/auth?sslmode=disable' up
create_migrate:
	migrate create -ext sql -dir ./db/migrations $(migrate_name)

build-image:
	docker build -t morozvol/auth_service:latest
start-container:
	docker run --name morozvol/auth_service:latest -p 8050:8050 --env-file .env


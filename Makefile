

MIGRATIONS_PATH := migrations
POSTGRES_DB := postgres://postgres:1234@localhost:5432/chatbox?sslmode=disable

gen:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	./proto/chat.proto

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migfile:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq chat

up:
	migrate -path $(MIGRATIONS_PATH) -database $(POSTGRES_DB) up

down:
	migrate -path $(MIGRATIONS_PATH) -database $(POSTGRES_DB) down

force:
	@if [ -z "$(version)" ]; then \
	  echo "Error: please specify version argument V"; \
	  exit 1; \
	fi
	migrate -path $(MIGRATIONS_PATH) -database $(POSTGRES_DB) force $(version)

.PHONY: gen install migfile up down force


\
make force version=1 \

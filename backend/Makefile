# Fill in these in your environment vars
#PSQL_DATABASE
#PSQL_USER
#PSQL_PASSWORD
#PSQL_HOST
#PSQL_PORT ?= 5432

PSQL_DSN ?= $(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(PGDATABASE)?sslmode=disable
PSQL_MIGRATION_DIR ?= internal/pkg/db/migrations/postgresql
# Version - this is optionally used on goto command
V?=
# Number of migrations - this is optionally used on up and down commands
N?=


local-db:
	@ docker-compose up -d


	@ until mysql --host=$(PSQL_HOST) --port=$(PSQL_PORT) --user=$(PSQL_USER) -p$(PSQL_PASSWORD) --protocol=tcp -e 'SELECT 1' >/dev/null 2>&1 && exit 0; do \
	  >&2 echo "MySQL is unavailable - sleeping"; \
	  sleep 5 ; \
	done


	@ echo "MySQL is up and running!"

docker-db:
	@ docker exec -it postgres psql -U $(PGUSER) -d $(PGDATABASE)

migrate-setup:
	@if [ -z "$$(which migrate)" ]; then echo "Installing migrate command..."; go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; fi

migrate-create:
	@ migrate create -ext sql -dir $(PSQL_MIGRATION_DIR) -seq $(Name)

migrate-up: migrate-setup
	@ migrate -database 'postgres://$(PSQL_DSN)' -path $(PSQL_MIGRATION_DIR) up $(N)


migrate-down: migrate-setup
	@ migrate -database 'postgres://$(PSQL_DSN)' -path $(PSQL_MIGRATION_DIR) down $(N)


migrate-to-version: migrate-setup
	@ migrate -database 'postgres://$(PSQL_DSN)' -path $(PSQL_MIGRATION_DIR) goto $(V)


drop-db: migrate-setup
	@ migrate -database 'postgres://$(PSQL_DSN)' -path $(PSQL_MIGRATION_DIR) drop


force-version: migrate-setup
	@ migrate -database 'postgres://$(PSQL_DSN)' -path $(PSQL_MIGRATION_DIR) force $(V)


migration-version: migrate-setup
	@ migrate -database 'postgres://$(PSQL_DSN)' -path $(PSQL_MIGRATION_DIR) version


gqlgen-gen:
	go run github.com/99designs/gqlgen@latest generate


gqlgen-init:
	go run github.com/99designs/gqlgen@latest init


build:server.go
	@ go build -o build/ .

dev:server.go
	@ go run ./server.go

run: build
	@ ./build/hackernews-go-graphql

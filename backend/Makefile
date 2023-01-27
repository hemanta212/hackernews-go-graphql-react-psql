# Fill in these in your environment vars
#MYSQL_DATABASE
#MYSQL_USER
#MYSQL_PASSWORD
#MYSQL_HOST
#MYSQL_PORT ?= 5432

MYSQL_DSN ?= $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)
MYSQL_MIGRATION_DIR ?= internal/pkg/db/migrations/mysql
# Version - this is optionally used on goto command
V?=
# Number of migrations - this is optionally used on up and down commands
N?=


local-db:
	@ docker-compose up -d


	@ until mysql --host=$(MYSQL_HOST) --port=$(MYSQL_PORT) --user=$(MYSQL_USER) -p$(MYSQL_PASSWORD) --protocol=tcp -e 'SELECT 1' >/dev/null 2>&1 && exit 0; do \
	  >&2 echo "MySQL is unavailable - sleeping"; \
	  sleep 5 ; \
	done


	@ echo "MySQL is up and running!"


migrate-setup:
	@if [ -z "$$(which migrate)" ]; then echo "Installing migrate command..."; go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate; fi

migrate-create:
	@ migrate create -ext sql -dir $(MYSQL_MIGRATION_DIR) -seq $(Name)

migrate-up: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATION_DIR) up $(N)


migrate-down: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATION_DIR) down $(N)


migrate-to-version: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATION_DIR) goto $(V)


drop-db: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATION_DIR) drop


force-version: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATION_DIR) force $(V)


migration-version: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATION_DIR) version


gqlgen-gen:
	go run github.com/99designs/gqlgen@latest generate


gqlgen-init:
	go run github.com/99designs/gqlgen@latest init


build:
	@ go build -o build/ .


run: build
	@ ./build/hackernews-go-graphql

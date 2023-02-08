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

dev:
	@ go run ./server.go

run: migrate-up
	@ go build -o ./build/ .
	@ ./build/hackernews-go-graphql

docker-run: vendor run

vendor:
	go mod tidy
	go mod download
	go mod vendor
	go mod verify

gqlgen-gen:
	go run github.com/99designs/gqlgen@latest generate

gqlgen-init:
	go run github.com/99designs/gqlgen@latest init
	
docker-up: vendor
	docker build -t hackernews-dev .
	docker run -it --rm -v $(HTTPS_SSL)/fullchain.pem:/etc/letsencrypt/live/vps.osac.org.np/fullchain.pem -v $(HTTPS_SSL)/privkey.pem:/etc/letsencrypt/live/vps.osac.org.np/privkey.pem -v $(PWD):/go/src/backend --net=host --env-file=.env hackernews-dev

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



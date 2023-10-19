PSQL_DSN ?= $(PGUSER):$(PGPASSWORD)@$(PGHOST):$(PGPORT)/$(PGDATABASE)?sslmode=disable
PSQL_MIGRATION_DIR ?= internal/pkg/db/migrations/postgresql
# Version - this is optionally used on goto command
V?=
# Number of migrations - this is optionally used on up and down commands
N?=

APP?=hackernews
PORT?=8080
RELEASE?=0.0.1


dev: migrate-up
	@ go run ./server.go

build: vendor
	@ go build -o ./build/ .

run: build migrate-up
	@ ./build/hackernews-go-graphql

vendor:
	go mod tidy
	go mod download
	go mod vendor
	go mod verify

gqlgen-gen:
	go run github.com/99designs/gqlgen@latest generate

gqlgen-init:
	go run github.com/99designs/gqlgen@latest init

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

container: 
	docker compose -f ../docker-compose.prod.yml build

minikube: container
	minikube image load $(APP):$(RELEASE)
	for t in $(shell find ../kubernetes/backend -type f -name "*.yaml"); do \
	cat $$t | \
    	sed -E "s/\{ \{(\s*)\.Release(\s*)\} \}/$(RELEASE)/g" | \
    	sed -E "s/\{ \{(\s*)\.ServiceName(\s*)\} \}/$(APP)/g"; \
	echo ---; \
    done > tmp.yaml

	minikube kubectl -- apply -f tmp.yaml
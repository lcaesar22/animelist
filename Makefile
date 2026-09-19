
#============================================================================
# HELPERS
#============================================================================

## help : print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##/p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'

#& create a new confirm target
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [Y/N]' && read ans && [ $${ans-N} = Y ]

#============================================================================
# DEVELOPMENT
#============================================================================

## run/api: run the cmd/api application
.PHONY: run
run:
	go run ./cmd/api -db-dsn=${ANIMELIST_DB_DSN}

## psql: connect to the database using psql
.PHONY: psql
psql:
	psql ${ANIMELIST_DB_DSN}

## db/migrations/new name = $1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Running up migration files for ${name}'
	migrate -path ./migrations -database ${name} up

## db/migrations/up: apply all database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations....'
	migrate -path ./migrations -database ${ANIMELIST_DB_DSN} up

#===============================================================================
# QUALITY CONTROL
#================================================================================

## tidy: tidy module dependencies for the go files
.PHONY: tidy
tidy:
	@echo 'Tidying up module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies..'
	go mod verify
	go mod vendor
	@echo 'Formatting .go files'
	go fmt ./...

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies...'
	go mod tidy -diff
	go mod verify
	@echo 'Verifying code ...'
	go vet ./...
	go tool staticcheck ./...
	@echo 'Running tests....'
	go test -race -vet=off ./...
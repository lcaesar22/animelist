# animelist
# animelist
# animelist api
A RESTful JSON API for managing an anime watchlist - built in Go
Users can register, activate their accounts, and manage their anime collection with full CRUD operations, filtering, sorting, pagination, token-based authentication with fine-grained permissions.

## Features
- RESTful endpoints versioned under `/v1`
- Full CRUD for anime records (title, year, genres, runtime, rating, version)
- Filtering, full-text search, sorting, and pagination with metadata
- User registration with email activation tokens
- Token-based authentication(`Bearer <token>`)
- Permission-based authorisation
- Rate limiting(global and per-IP)
- Structured JSON logging and panic recovery
- Graceful shutdown
- Application metrics exposed at `/debug/vars`
- CORS support
- SQL migrations for PostgreSQL

## Stack
- Go 1.25
- PostgreSQL
- Std library
- Chi router

## Getting Started
- Go 1.25+
- PostgeSQL
- `curl` (for the endpoints)
- `make` (optional, for Makefile targets)

Configurations
Flag                            Default
-port                           4000
-env                            development
-db-dsn                         -
-db-max-open-conns              25
-db-max-idle-conns              25
-db-max-idle-time               15m
-smtp-host                      -
-smtp-port                      -
-smtp-username                  -
-smtp-password                  -
-rate-limiter-enabled           true
-rate-limiter-rps               2
-rate-limiter-burst             4
-version                        1.0.0

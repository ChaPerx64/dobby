# About

This project is a little helper for a household that helps managing finances and expenses.

# Code structure

All of the project's services and their corresponding configurations are stored in `apps/` in corresponding subdirectories.

# Services

## Persistence (Database)

PostgreSQL (v17) is used for persistence.


### DB schema principles

- Schema should be at least 3NF normal.
- Database should generate NO data. Including, but not limited to:
  - No ID generation
  - No timestamp generation
  - No default values
- Database is used for persistence and atomicity, so its schema should align closely with Service-level models

### Migrations

Schema migrations are managed by DBMate, running in Docker as a part of `docker-compose`.
DBMate must not be installed as a local tooling/dependency, but be called in/from Docker container.

Migrations MUST be in SQL, be idempotent and contain up and down directions.

Migration files can be found in ./apps/db/migrations/

Filename should adhere to this format: `YYYY-MM-DD_HH-MM_<short-description>.sql`

For other information, consult: https://github.com/amacneil/dbmate/blob/main/README.md

### Deployment

PostgreSQL instance is run in a Docker container as a part of docker-compose environment.

For development, you must run it (or ensure it is running) with `docker compose up -d dobby-db`.

To apply migrations to this DB, run `docker compose run --rm dbmate`

## Back-end (`apps/backend/`)

Back-end is a REST API service for business logic and data management.
It is implemented with Golang and realises "Ports and Adapters" (aka "Hexagonal") architecture.

### Adapters/Repository layer

This layer must be implemented, utilizing "Repositories" pattern.
ORM use is explicitly avoided in favor of writing pure SQL queries.

Repositories must accept and return Service-layer models and DTOs.

Repository methods MUST NOT generate any data, meaning, they should have no default values.

Repository for DB interaction is defined in `./apps/backend/internal/adapters/persistence/`

Repositories must get UnitOfWork object from Context.

### Service layer

This is where most of calculation happens and data is transformed and generated.

This is also where is where data models are defined.
All other layers must try to use these models as much as prudently possible.

The service layer should control a UnitOfWork: begin transactions and end them.

### Adapters/API layer

API layer-related code can be found in `apps/backend/internal/api/`.
API layer is implemented with `ogen` auto-generated types and stubs for endpoints from the OpenAPI specification.
To regenerate, run `go generate ./...`

Generated stubs and types are located in `./apps/backend/internal/adapters/oas/`

Adapter layer should populate service-level models, where possible, with models themselves defined in the service layer.
In case where an API-layer DTO must differ from Service-layer (due to API differences, e.g. added defaults or deprecated fields),
it must define `toLogicModel` and `fromLogicModel` methods which then will be called by the API layer.

### Exceptions

Exceptions must be defined at Service layer, with following rules for integration:

- Driven adapters must use Service-layer exceptions to wrap inner exceptions.
- Driver adapters must wrap Service-layer exceptions into appropriate exceptions.

Service layer exceptions MUST NOT "know" of their driver-adapter corresponding errors.
E.g. if Service/Domain-layer defines AllocationNotFound exception, it should NOT know of HTTP code,
instead, an API adapter should discern and wrap it into a 404 Not Found error.


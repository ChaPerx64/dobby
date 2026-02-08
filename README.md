# About

This project is a little helper for a household that helps managing finances and expenses.

# Code structure

All of the project's services and their corresponding configurations are stored in `apps/`

# Services

## Back-end (`apps/backend/`)

Back-end is a REST API service for business logic and data management.
It is implemented with Golang and realises "Adapters" (aka "Hexagonal") architecture.

### API layer

API layer-related code can be found in `apps/backend/internal/api/`.
API layer is implemented using `ogen` that auto-generates types and stubs for endpoints from the OpenAPI specification.
To regenerate, run `go generate ./...`




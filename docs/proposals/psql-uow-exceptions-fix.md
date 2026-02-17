# Proposal: PostgreSQL Adapter, Unit of Work, and Error Handling Refinement

This proposal outlines the implementation details for bringing the Dobby backend into full compliance with the architecture defined in the project `README.md`.

## 1. Unit of Work (UoW) Pattern

To satisfy the requirement that "Repositories must get UnitOfWork object from Context", we will implement a context-based transaction management system.

### Context Key and Type
```go
type uowKey struct{}

// UnitOfWork defines the interface for managing transactions.
type UnitOfWork interface {
	Commit() error
	Rollback() error
}
```

### Repository Integration
The `Repository` interface will be updated (or its implementation will be designed) to extract the active transaction (or DB connection) from the context.

```go
func (r *psqlRepo) getDB(ctx context.Context) DB {
	if tx, ok := ctx.Value(uowKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.db
}
```

### Service Layer Control
The service layer will use a `TransactionManager` to wrap operations:

```go
func (s *dobbyFinancier) RecordTransaction(ctx context.Context, t Transaction) (*Transaction, error) {
	err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		// All repo calls here use the tx in ctx
		return s.repo.SaveTransaction(ctx, &t)
	})
	return &t, err
}
```

## 2. PostgreSQL Adapter implementation

A new adapter will be created in `apps/backend/internal/adapters/persistence/postgres.go`.

### Key Features:
- **Pure SQL**: No ORM, as per README mandates.
- **UUID Support**: Using `github.com/google/uuid`.
- **Cents for Currency**: `Amount` stored as `BIGINT`.
- **No Data Generation**: IDs and timestamps will be passed from the service layer.

### Schema Alignment:
The adapter will target the schema defined in `apps/db/schema.sql`.

## 3. Exception/Error Handling Refinement

We will implement a centralized error mapper in the API layer to translate Service-layer exceptions into HTTP-aware `oas.ErrorStatusCode`.

### Mapping Table:
| Service Error | HTTP Status |
| :--- | :--- |
| `ErrNotFound` | 404 Not Found |
| `ErrValidation` | 400 Bad Request |
| `ErrPeriodOverlap` | 409 Conflict |
| `ErrInsufficientFunds` | 422 Unprocessable Entity |
| *Default* | 500 Internal Server Error |

### Implementation in `handlers.go`:
```go
func (h *dobbyHandler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	var code int
	switch {
	case errors.Is(err, service.ErrNotFound):
		code = 404
	case errors.Is(err, service.ErrValidation):
		code = 400
	// ... etc
	default:
		code = 500
	}
	return &oas.ErrorStatusCode{
		StatusCode: code,
		Response: oas.Error{
			Code:    int64(code),
			Message: err.Error(),
		},
	}
}
```

## 4. DTO Mapping

Following the "toLogicModel/fromLogicModel" suggestion, we will implement these as methods on the auto-generated OAS types. 

To prevent these methods from being overwritten during code regeneration:
- They will be defined in a new, non-generated file: `apps/backend/internal/adapters/oas/mapping.go`.
- This file will reside in the same `oas` package as the generated code, allowing us to attach methods to the DTO structs.

Example:
```go
// apps/backend/internal/adapters/oas/mapping.go
package oas

func (req *CreateEnvelope) ToLogicModel() string {
    return req.Name
}
```

## 5. Implementation Steps

1. **Step 1: PostgreSQL Infrastructure**: Set up the DB connection pool and the new `postgres.go` repository.
2. **Step 2: Transaction Manager**: Implement the UoW/Transaction logic in `internal/adapters/persistence`.
3. **Step 3: Service Layer Refactor**: Update `dobbyFinancier` to use the transaction manager for multi-step operations.
4. **Step 4: Error Handling**: Update `dobbyHandler.NewError` and refine all API endpoints to return correct error codes.
5. **Step 5: Wiring**: Update `main.go` to use the PostgreSQL repository instead of the memory repository.

# Service Layer Implementation Proposal

This document outlines the proposed structure and models for the Service Layer of the Dobby backend, adhering to the "Ports and Adapters" architecture described in the project README.

## 1. Domain Models

These models represent the core business entities and will be defined in the service layer (`apps/backend/internal/service`). They are pure data structures decoupled from database schema details (like ORM tags) and API transport details (like JSON tags).

### Entities (Persisted)

```go
package service

import (
	"time"
	"github.com/google/uuid"
)

// User represents a household member.
type User struct {
	ID   uuid.UUID
	Name string
}

// Period represents a defined financial timeframe.
type Period struct {
	ID        uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

// Envelope represents a budget category/bucket (e.g., "Groceries").
type Envelope struct {
	ID   uuid.UUID
	Name string
}

// Transaction represents a financial movement.
type Transaction struct {
	ID          uuid.UUID
	PeriodID    uuid.UUID
	EnvelopeID  uuid.UUID
	Amount      int64     // Stored in cents. Positive = Income (Budget), Negative = Expense.
	Description string
	Date        time.Time
	Category    string    // Analytics tag
}
```

### Aggregates & Value Objects (Calculated)

These structures are used to pass rich data back to the API layer without coupling the internal calculation logic to the API response format.

```go
// PeriodSummary enriches the Period entity with calculated financial status.
type PeriodSummary struct {
	Period                 Period
	TotalBudget            int64 // Sum of all positive transactions (Income) across all envelopes
	TotalSpent             int64 // Sum of all negative transactions (Expense) across all envelopes
	TotalRemaining         int64 // TotalBudget + TotalSpent (since spent is negative)
	ProjectedEndingBalance int64 // Forecast logic
	EnvelopeStats          []EnvelopeStat
}

// EnvelopeStat provides a snapshot of an envelope's performance within a specific period.
type EnvelopeStat struct {
	Envelope  Envelope
	Allocated int64 // Sum of positive transactions (Income) for this envelope in this period
	Spent     int64 // Sum of negative transactions (Expense) for this envelope in this period
	Remaining int64 // Allocated + Spent (since spent is negative)
}
```

## 2. Service Interface

The Service definition acts as the "Inbound Port" (or primary port) for the application core.

```go
package service

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type FinanceService interface {
	// Period Operations
	CreatePeriod(ctx context.Context, start, end time.Time) (*Period, error)
	GetPeriodSummary(ctx context.Context, id uuid.UUID) (*PeriodSummary, error)
	ListPeriods(ctx context.Context) ([]Period, error)

	// Transaction Operations
	RecordTransaction(ctx context.Context, t Transaction) (*Transaction, error)
	
	// Envelope Operations
	CreateEnvelope(ctx context.Context, name string) (*Envelope, error)
	ListEnvelopes(ctx context.Context) ([]Envelope, error)
}
```

## 3. Repository Interface (Outbound Ports)

Repositories act as "Outbound Ports". As per the README, they must utilize a `UnitOfWork` from the context and return Service-layer models.

```go
type Repository interface {
	// Transactional support
	WithTx(ctx context.Context, fn func(repo Repository) error) error

	// Domain methods
	SaveUser(ctx context.Context, u *User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	
	SavePeriod(ctx context.Context, p *Period) error
	GetPeriod(ctx context.Context, id uuid.UUID) (*Period, error)
	ListPeriods(ctx context.Context) ([]Period, error)
	
	SaveEnvelope(ctx context.Context, e *Envelope) error
	ListEnvelopes(ctx context.Context) ([]Envelope, error)
	
	SaveTransaction(ctx context.Context, t *Transaction) error
	ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
}
```

## 4. Error Handling

Service-level exceptions must be defined here. They should be agnostic of the transport layer (no HTTP codes).

```go
var (
	ErrNotFound          = errors.New("resource not found")
	ErrValidation        = errors.New("validation error")
	ErrPeriodOverlap     = errors.New("period dates overlap with existing period")
	ErrInsufficientFunds = errors.New("insufficient funds") // If applicable logic exists
)
```

## 5. Implementation Strategy

1.  **Define Models**: Define the structs in `apps/backend/internal/service/models.go`.
2.  **Define Interfaces**: Update `apps/backend/internal/service/service.go` (or similar) with `FinanceService` and `Repository` interfaces.
3.  **Implement Service**: Update `apps/backend/internal/service/dobbyFinancier.go` (or rename to `service.go` implementation) to implement `FinanceService`. This will contain the business logic (e.g., calculating `TotalBudget` by summing positive transactions).
4.  **Wiring**: The API adapter (ogen handlers) will call `FinanceService`. The `FinanceService` will call the `Repository`.

package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrValidation        = errors.New("validation error")
	ErrPeriodOverlap     = errors.New("period dates overlap with existing period")
	ErrInsufficientFunds = errors.New("insufficient funds")
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

type TransactionFilter struct {
	PeriodID   *uuid.UUID
	EnvelopeID *uuid.UUID
}

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

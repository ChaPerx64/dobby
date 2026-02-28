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
	ErrConflict          = errors.New("resource conflict")
)

type FinanceService interface {
	// Period Operations
	CreatePeriod(ctx context.Context, start, end *time.Time) (*Period, error)
	GetCurrentPeriod(ctx context.Context) (*PeriodSummary, error)
	GetPeriodSummary(ctx context.Context, id uuid.UUID) (*PeriodSummary, error)
	ListPeriods(ctx context.Context) ([]Period, error)

	// Transaction Operations
	RecordTransaction(ctx context.Context, t Transaction) (*Transaction, error)
	ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*Transaction, error)
	UpdateTransaction(ctx context.Context, t Transaction) (*Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error

	// Envelope Operations
	CreateEnvelope(ctx context.Context, name string) (*Envelope, error)
	ListEnvelopes(ctx context.Context) ([]Envelope, error)
	DeleteEnvelope(ctx context.Context, id uuid.UUID) error
}

type TransactionFilter struct {
	PeriodID   *uuid.UUID
	EnvelopeID *uuid.UUID
}

type TransactionManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository interface {
	// Domain methods
	SaveUser(ctx context.Context, u *User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)

	SavePeriod(ctx context.Context, p *Period) error
	GetPeriod(ctx context.Context, id uuid.UUID) (*Period, error)
	GetCurrentPeriod(ctx context.Context) (*Period, error)
	ListPeriods(ctx context.Context) ([]Period, error)

	SaveEnvelope(ctx context.Context, e *Envelope) error
	ListEnvelopes(ctx context.Context) ([]Envelope, error)
	DeleteEnvelope(ctx context.Context, id uuid.UUID) error

	SaveTransaction(ctx context.Context, t *Transaction) error
	ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) error

	GetPeriodStats(ctx context.Context, periodID uuid.UUID) ([]EnvelopeStat, error)
}

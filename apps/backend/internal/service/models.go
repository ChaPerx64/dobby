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
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
}

// Transaction represents a financial movement.
type Transaction struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	PeriodID    uuid.UUID
	EnvelopeID  uuid.UUID
	Amount      int64 // Stored in cents. Positive = Income (Budget), Negative = Expense.
	Description string
	Date        time.Time
	Category    string // Analytics tag
}

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

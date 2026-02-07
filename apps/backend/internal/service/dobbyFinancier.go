package service

import (
	"time"

	"github.com/google/uuid"
)

type EnvelopeSummary struct {
	UserID     uuid.UUID
	EnvelopeID uuid.UUID
	Envelope   string
	Amount     int64
	Spent      int64
	Remaining  int64
}

type FinancialPeriod struct {
	ID        uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool

	TotalBudget            int64
	TotalSpent             int64
	TotalRemaining         int64
	ProjectedEndingBalance int64
	EnvelopeSummaries      []EnvelopeSummary
}

type PeriodListItem struct {
	ID        uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
}

type DobbyFinancierService struct{}

func (dfs DobbyFinancierService) GetCurrentPeriod() (FinancialPeriod, error) {
	return FinancialPeriod{}, nil
}

func (dfs DobbyFinancierService) ListPeriods() ([]PeriodListItem, error) {
	return []PeriodListItem{}, nil
}

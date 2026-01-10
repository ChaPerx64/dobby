package service

import (
	"time"

	"github.com/google/uuid"
)

type FundCategory struct {
	ID             uuid.UUID
	AllocationUser uuid.UUID
	Purpose        string
	FundsAllocated int
	FundsSpent     int
	FundsRemaining int
}

type FinancialPeriod struct {
	ID        uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool

	TotalBudgetCents int
	Allocations      []FundCategory
}

type DobbyFinancierService struct{}

func (dfs DobbyFinancierService) GetCurrentPeriod() (FinancialPeriod, error) {
	return FinancialPeriod{}, nil
}

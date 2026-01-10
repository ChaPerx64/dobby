package api

import (
	"context"
	"log"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
	"github.com/ChaPerx64/dobby/apps/backend/internal/service"
	"github.com/google/uuid"
)

func (h dobbyHandler) GetCurrentUser(ctx context.Context) (r *oas.User, _ error) {
	log.Println("Got a request @/me")
	return &oas.User{
		ID:             uuid.UUID{},
		Name:           "mock name",
		CurrentBalance: 10000,
	}, nil
}

func (h dobbyHandler) GetCurrentPeriod(ctx context.Context) (r *oas.Period, _ error) {
	log.Println("Got a request @/me/period")
	financialPeriod, err := service.DobbyFinancierService{}.GetCurrentPeriod()
	if err != nil {
		return nil, err
	}

	// Convert service.FinancialPeriod to oas.Period
	allocations := make([]oas.Allocation, len(financialPeriod.Allocations))
	for i, a := range financialPeriod.Allocations {
		allocations[i] = oas.Allocation{
			ID:        a.ID,
			UserId:    a.AllocationUser,
			PeriodId:  financialPeriod.ID,
			Purpose:   oas.AllocationPurpose(a.Purpose),
			Amount:    int64(a.FundsAllocated),
			Spent:     int64(a.FundsSpent),
			Remaining: int64(a.FundsRemaining),
		}
	}

	return &oas.Period{
		ID:          financialPeriod.ID,
		StartDate:   financialPeriod.StartDate,
		EndDate:     financialPeriod.EndDate,
		IsActive:    financialPeriod.IsActive,
		TotalBudget: int64(financialPeriod.TotalBudgetCents),
		Allocations: allocations,
	}, nil
}

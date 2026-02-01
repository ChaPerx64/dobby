package api

import (
	"context"
	"log"
	"time"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
	"github.com/google/uuid"
)

// Ensure dobbyHandler implements oas.Handler
var _ oas.Handler = (*dobbyHandler)(nil)

func (h dobbyHandler) GetCurrentUser(ctx context.Context) (*oas.User, error) {
	log.Println("Got a request @/me")
	return &oas.User{
		ID:             uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Name:           "TheMan",
		CurrentBalance: 10000,
	}, nil
}

func (h dobbyHandler) GetCurrentPeriod(ctx context.Context) (*oas.Period, error) {
	log.Println("Got a request @/periods/current")
	// For now, mock the response since service layer is not fully ready for the new schema
	// financialPeriod, err := service.DobbyFinancierService{}.GetCurrentPeriod()

	return &oas.Period{
		ID:                     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		StartDate:              time.Now().AddDate(0, 0, -5),
		EndDate:                time.Now().AddDate(0, 1, -5),
		IsActive:               true,
		TotalBudget:            12000000,
		TotalSpent:             2300000,
		TotalRemaining:         9700000,
		ProjectedEndingBalance: oas.NewOptInt64(-300000), // Optional field
	}, nil
}

func (h dobbyHandler) ListAllocations(ctx context.Context, params oas.ListAllocationsParams) ([]oas.Allocation, error) {
	log.Println("Got a request @/allocations")
	
	// Mock Envelopes
	groceriesID := uuid.MustParse("22222222-2222-2222-2222-222222222201")
	pocketChaianID := uuid.MustParse("22222222-2222-2222-2222-222222222202")
	
	periodID := params.PeriodId
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if params.UserId.IsSet() {
		userID = params.UserId.Value
	}

	return []oas.Allocation{
		{
			ID:         uuid.New(),
			PeriodId:   periodID,
			UserId:     userID,
			EnvelopeId: groceriesID,
			Envelope:   "Groceries",
			Amount:     5000000,
			Spent:      1200000,
			Remaining:  3800000,
		},
		{
			ID:         uuid.New(),
			PeriodId:   periodID,
			UserId:     userID,
			EnvelopeId: pocketChaianID,
			Envelope:   "Chaian's pocket money",
			Amount:     3500000,
			Spent:      800000,
			Remaining:  2700000,
		},
	}, nil
}

func (h dobbyHandler) ListEnvelopes(ctx context.Context) ([]oas.Envelope, error) {
	log.Println("Got a request @/envelopes")
	return []oas.Envelope{
		{
			ID:   uuid.MustParse("22222222-2222-2222-2222-222222222201"),
			Name: "Groceries",
		},
		{
			ID:   uuid.MustParse("22222222-2222-2222-2222-222222222202"),
			Name: "Chaian's pocket money",
		},
		{
			ID:   uuid.MustParse("22222222-2222-2222-2222-222222222203"),
			Name: "Sophia's pocket money",
		},
	}, nil
}

// Stub other methods (optional, as UnimplementedHandler handles them, but explicit logging is nice)
func (h dobbyHandler) ListPeriods(ctx context.Context) ([]oas.Period, error) {
	log.Println("Got a request @/periods")
	p, _ := h.GetCurrentPeriod(ctx)
	return []oas.Period{*p}, nil
}

// Implement other required methods by delegating to UnimplementedHandler (implicit)
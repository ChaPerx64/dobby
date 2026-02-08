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

	idStr, ok := GetUserID(ctx)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userName := "TheMan"

	if ok {
		log.Printf("Authenticated User ID: %s\n", idStr)
		if parsed, err := uuid.Parse(idStr); err == nil {
			userID = parsed
		}
		userName = "Authenticated User"
	}

	return &oas.User{
		ID:             userID,
		Name:           userName,
		CurrentBalance: 10000,
	}, nil
}

func (h dobbyHandler) GetCurrentPeriod(ctx context.Context) (*oas.Period, error) {
	log.Println("Got a request @/periods/current")
	
	groceriesID := uuid.MustParse("22222222-2222-2222-2222-222222222201")
	pocketChaianID := uuid.MustParse("22222222-2222-2222-2222-222222222202")
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	return &oas.Period{
		ID:             uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		StartDate:      time.Now().AddDate(0, 0, -5),
		EndDate:        time.Now().AddDate(0, 1, -5),
		IsActive:       true,
		TotalBudget:    12000000,
		TotalSpent:     2300000,
		TotalRemaining: 9700000,
		ProjectedEndingBalance: oas.NewOptInt64(-300000),
		EnvelopeSummaries: []oas.EnvelopeSummary{
			{
				UserId:     userID,
				EnvelopeId: groceriesID,
				Envelope:   "Groceries",
				Amount:     5000000,
				Spent:      1200000,
				Remaining:  3800000,
			},
			{
				UserId:     userID,
				EnvelopeId: pocketChaianID,
				Envelope:   "Chaian's pocket money",
				Amount:     3500000,
				Spent:      800000,
				Remaining:  2700000,
			},
		},
	}, nil
}

func (h dobbyHandler) GetPeriod(ctx context.Context, params oas.GetPeriodParams) (oas.GetPeriodRes, error) {
	log.Printf("Got a request @/periods/%s\n", params.PeriodId)
	
	// For now, return the current period mock for any ID
	p, _ := h.GetCurrentPeriod(ctx)
	p.ID = params.PeriodId
	return p, nil
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

func (h dobbyHandler) ListPeriods(ctx context.Context) ([]oas.PeriodListItem, error) {
	log.Println("Got a request @/periods")
	p, _ := h.GetCurrentPeriod(ctx)
	return []oas.PeriodListItem{
		{
			ID:        p.ID,
			StartDate: p.StartDate,
			EndDate:   p.EndDate,
			IsActive:  p.IsActive,
		},
	}, nil
}

func (h dobbyHandler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	return &oas.ErrorStatusCode{
		StatusCode: 500,
		Response: oas.Error{
			Code:    500,
			Message: err.Error(),
		},
	}
}

package api

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/oas"
	"github.com/ChaPerx64/dobby/apps/backend/internal/service"
	"github.com/google/uuid"
)

// Ensure dobbyHandler implements oas.Handler
var _ oas.Handler = (*dobbyHandler)(nil)

func (h *dobbyHandler) GetCurrentUser(ctx context.Context) (*oas.User, error) {
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
		ID:   userID,
		Name: userName,
	}, nil
}

func (h *dobbyHandler) GetCurrentPeriod(ctx context.Context) (*oas.Period, error) {
	log.Println("Got a request @/periods/current")

	// For now, if no periods exist, return a mock or empty
	periods, err := h.financeService.ListPeriods(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	if len(periods) == 0 {
		// Mock a period for now if none exists to maintain original behavior
		return &oas.Period{
			ID:                     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			StartDate:              time.Now().AddDate(0, 0, -5),
			EndDate:                time.Now().AddDate(0, 1, -5),
			TotalBudget:            12000000,
			TotalSpent:             2300000,
			TotalRemaining:         9700000,
			ProjectedEndingBalance: oas.NewOptInt64(-300000),
		}, nil
	}

	// Assuming the last one is "current" for simplicity of mock refactor
	summary, err := h.financeService.GetPeriodSummary(ctx, periods[len(periods)-1].ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	return mapPeriodSummaryToOAS(summary), nil
}

func (h *dobbyHandler) GetPeriod(ctx context.Context, params oas.GetPeriodParams) (oas.GetPeriodRes, error) {
	log.Printf("Got a request @/periods/%s\n", params.PeriodId)

	summary, err := h.financeService.GetPeriodSummary(ctx, params.PeriodId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return &oas.GetPeriodNotFound{}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return mapPeriodSummaryToOAS(summary), nil
}

func (h *dobbyHandler) ListEnvelopes(ctx context.Context) ([]oas.Envelope, error) {
	log.Println("Got a request @/envelopes")
	envelopes, err := h.financeService.ListEnvelopes(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	res := make([]oas.Envelope, len(envelopes))
	for i, e := range envelopes {
		res[i] = oas.Envelope{
			ID:   e.ID,
			Name: e.Name,
		}
	}
	return res, nil
}

func (h *dobbyHandler) ListPeriods(ctx context.Context) ([]oas.PeriodListItem, error) {
	log.Println("Got a request @/periods")
	periods, err := h.financeService.ListPeriods(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	res := make([]oas.PeriodListItem, len(periods))
	for i, p := range periods {
		res[i] = oas.PeriodListItem{
			ID:        p.ID,
			StartDate: p.StartDate,
			EndDate:   p.EndDate,
		}
	}
	return res, nil
}

func (h *dobbyHandler) CreateEnvelope(ctx context.Context, req *oas.CreateEnvelope) (*oas.Envelope, error) {
	log.Println("Got a request POST /envelopes")
	env, err := h.financeService.CreateEnvelope(ctx, req.Name)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &oas.Envelope{
		ID:   env.ID,
		Name: env.Name,
	}, nil
}

func (h *dobbyHandler) CreatePeriod(ctx context.Context, req *oas.CreatePeriod) (*oas.Period, error) {
	log.Println("Got a request POST /periods")
	p, err := h.financeService.CreatePeriod(ctx, req.StartDate, req.EndDate)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	summary, err := h.financeService.GetPeriodSummary(ctx, p.ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	return mapPeriodSummaryToOAS(summary), nil
}

func (h *dobbyHandler) CreateTransaction(ctx context.Context, req *oas.CreateTransaction) (oas.CreateTransactionRes, error) {
	log.Println("Got a request POST /transactions")

	// We need to find which period this transaction belongs to.
	// For now, let's assume it belongs to the period covering its date.
	// Or just use the current period if we don't have a lookup.
	// The service.Transaction model needs a PeriodID.

	tDate := time.Now()
	if v, ok := req.Date.Get(); ok {
		tDate = v
	}

	periods, err := h.financeService.ListPeriods(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	var periodID uuid.UUID
	for _, p := range periods {
		if (tDate.After(p.StartDate) || tDate.Equal(p.StartDate)) && (tDate.Before(p.EndDate) || tDate.Equal(p.EndDate)) {
			periodID = p.ID
			break
		}
	}

	if periodID == uuid.Nil {
		return nil, h.NewError(ctx, errors.New("no period found for transaction date"))
	}

	t, err := h.financeService.RecordTransaction(ctx, service.Transaction{
		PeriodID:    periodID,
		EnvelopeID:  req.EnvelopeId,
		Amount:      req.Amount,
		Description: req.Description.Or(""),
		Date:        tDate,
		Category:    req.Category.Or(""),
	})
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	return &oas.Transaction{
		ID:          t.ID,
		PeriodId:    t.PeriodID,
		EnvelopeId:  t.EnvelopeID,
		Amount:      t.Amount,
		Description: oas.NewOptString(t.Description),
		Date:        t.Date,
		Category:    oas.NewOptString(t.Category),
	}, nil
}

func mapPeriodSummaryToOAS(s *service.PeriodSummary) *oas.Period {
	envSummaries := make([]oas.EnvelopeSummary, len(s.EnvelopeStats))
	for i, stat := range s.EnvelopeStats {
		envSummaries[i] = oas.EnvelopeSummary{
			EnvelopeId:   stat.Envelope.ID,
			EnvelopeName: stat.Envelope.Name,
			Amount:       stat.Allocated,
			Spent:        stat.Spent,
			Remaining:    stat.Remaining,
		}
	}

	return &oas.Period{
		ID:                     s.Period.ID,
		StartDate:              s.Period.StartDate,
		EndDate:                s.Period.EndDate,
		TotalBudget:            s.TotalBudget,
		TotalRemaining:         s.TotalRemaining,
		TotalSpent:             s.TotalSpent,
		ProjectedEndingBalance: oas.NewOptInt64(s.ProjectedEndingBalance),
		EnvelopeSummaries:      envSummaries,
	}
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

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

func (h *dobbyHandler) GetCurrentPeriod(ctx context.Context) (*oas.PeriodSummary, error) {
	log.Println("Got a request @/periods/current")

	p, err := h.financeService.GetCurrentPeriod(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	return mapPeriodSummaryToOAS(p), nil
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

func (h *dobbyHandler) getUserID(ctx context.Context) uuid.UUID {
	idStr, ok := GetUserID(ctx)
	if !ok {
		// Default user for now if not authenticated, or we could return error
		return uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}
	parsed, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.MustParse("00000000-0000-0000-0000-000000000001")
	}
	return parsed
}

func (h *dobbyHandler) CreateEnvelope(ctx context.Context, req *oas.CreateEnvelope) (*oas.Envelope, error) {
	log.Println("Got a request POST /envelopes")
	userID := h.getUserID(ctx)
	env, err := h.financeService.CreateEnvelope(ctx, userID, req.ToLogicModel())
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &oas.Envelope{
		ID:   env.ID,
		Name: env.Name,
	}, nil
}

func (h *dobbyHandler) CreatePeriod(ctx context.Context, req *oas.CreatePeriod) (*oas.PeriodSummary, error) {
	log.Println("Got a request POST /periods")
	p, err := h.financeService.CreatePeriod(ctx, &req.StartDate, &req.EndDate)
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

	t := req.ToLogicModel()
	if t.Date.IsZero() {
		t.Date = time.Now()
	}

	periods, err := h.financeService.ListPeriods(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	var periodID uuid.UUID
	for _, p := range periods {
		if (t.Date.After(p.StartDate) || t.Date.Equal(p.StartDate)) && (t.Date.Before(p.EndDate) || t.Date.Equal(p.EndDate)) {
			periodID = p.ID
			break
		}
	}

	if periodID == uuid.Nil {
		return nil, h.NewError(ctx, errors.New("no period found for transaction date"))
	}

	t.PeriodID = periodID
	userID := h.getUserID(ctx)
	recorded, err := h.financeService.RecordTransaction(ctx, userID, t)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	return &oas.Transaction{
		ID:          recorded.ID,
		PeriodId:    recorded.PeriodID,
		EnvelopeId:  recorded.EnvelopeID,
		Amount:      recorded.Amount,
		Description: oas.NewOptString(recorded.Description),
		Date:        recorded.Date,
		Category:    oas.NewOptString(recorded.Category),
	}, nil
}

func mapPeriodSummaryToOAS(s *service.PeriodSummary) *oas.PeriodSummary {
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

	return &oas.PeriodSummary{
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

func (h *dobbyHandler) NewError(ctx context.Context, err error) *oas.ErrorStatusCode {
	var code int
	switch {
	case errors.Is(err, service.ErrNotFound):
		code = 404
	case errors.Is(err, service.ErrValidation):
		code = 400
	case errors.Is(err, service.ErrPeriodOverlap):
		code = 409
	case errors.Is(err, service.ErrInsufficientFunds):
		code = 422
	default:
		code = 500
	}
	return &oas.ErrorStatusCode{
		StatusCode: code,
		Response: oas.Error{
			Code:    code,
			Message: err.Error(),
		},
	}
}

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

func (h *dobbyHandler) CreateEnvelope(ctx context.Context, req *oas.CreateEnvelope) (*oas.Envelope, error) {
	log.Println("Got a request POST /envelopes")
	env, err := h.financeService.CreateEnvelope(ctx, req.ToLogicModel())
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &oas.Envelope{
		ID:   env.ID,
		Name: env.Name,
	}, nil
}

func (h *dobbyHandler) DeleteEnvelope(ctx context.Context, params oas.DeleteEnvelopeParams) (oas.DeleteEnvelopeRes, error) {
	log.Printf("Got a request DELETE /envelopes/%s\n", params.EnvelopeId)

	err := h.financeService.DeleteEnvelope(ctx, params.EnvelopeId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return &oas.DeleteEnvelopeNotFound{}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return &oas.DeleteEnvelopeNoContent{}, nil
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
	recorded, err := h.financeService.RecordTransaction(ctx, t)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	return mapTransactionToOAS(recorded), nil
}

func (h *dobbyHandler) ListTransactions(ctx context.Context, params oas.ListTransactionsParams) ([]oas.Transaction, error) {
	log.Println("Got a request GET /transactions")

	filter := service.TransactionFilter{}
	if v, ok := params.PeriodId.Get(); ok {
		filter.PeriodID = &v
	}

	transactions, err := h.financeService.ListTransactions(ctx, filter)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	res := make([]oas.Transaction, len(transactions))
	for i, t := range transactions {
		res[i] = *mapTransactionToOAS(&t)
	}
	return res, nil
}

func (h *dobbyHandler) GetTransaction(ctx context.Context, params oas.GetTransactionParams) (oas.GetTransactionRes, error) {
	log.Printf("Got a request GET /transactions/%s\n", params.TransactionId)

	t, err := h.financeService.GetTransaction(ctx, params.TransactionId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return &oas.GetTransactionNotFound{}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return mapTransactionToOAS(t), nil
}

func (h *dobbyHandler) UpdateTransaction(ctx context.Context, req *oas.UpdateTransaction, params oas.UpdateTransactionParams) (oas.UpdateTransactionRes, error) {
	log.Printf("Got a request PATCH /transactions/%s\n", params.TransactionId)

	existing, err := h.financeService.GetTransaction(ctx, params.TransactionId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return &oas.UpdateTransactionNotFound{}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	req.ApplyToModel(existing)

	// If date changed, we might need to update PeriodID
	if _, ok := req.Date.Get(); ok {
		periods, err := h.financeService.ListPeriods(ctx)
		if err != nil {
			return nil, h.NewError(ctx, err)
		}

		var periodID uuid.UUID
		for _, p := range periods {
			if (existing.Date.After(p.StartDate) || existing.Date.Equal(p.StartDate)) && (existing.Date.Before(p.EndDate) || existing.Date.Equal(p.EndDate)) {
				periodID = p.ID
				break
			}
		}

		if periodID != uuid.Nil {
			existing.PeriodID = periodID
		}
	}

	updated, err := h.financeService.UpdateTransaction(ctx, *existing)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	return mapTransactionToOAS(updated), nil
}

func (h *dobbyHandler) DeleteTransaction(ctx context.Context, params oas.DeleteTransactionParams) (oas.DeleteTransactionRes, error) {
	log.Printf("Got a request DELETE /transactions/%s\n", params.TransactionId)

	err := h.financeService.DeleteTransaction(ctx, params.TransactionId)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return &oas.DeleteTransactionNotFound{}, nil
		}
		return nil, h.NewError(ctx, err)
	}

	return &oas.DeleteTransactionNoContent{}, nil
}

func mapTransactionToOAS(t *service.Transaction) *oas.Transaction {
	return &oas.Transaction{
		ID:          t.ID,
		PeriodId:    t.PeriodID,
		EnvelopeId:  t.EnvelopeID,
		Amount:      t.Amount,
		Description: oas.NewOptString(t.Description),
		Date:        t.Date,
		Category:    oas.NewOptString(t.Category),
	}
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
	case errors.Is(err, service.ErrPeriodOverlap), errors.Is(err, service.ErrConflict):
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

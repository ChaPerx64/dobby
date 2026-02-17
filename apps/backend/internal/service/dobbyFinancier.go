package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type dobbyFinancier struct {
	repo Repository
}

func NewDobbyFinancier(repo Repository) FinanceService {
	return &dobbyFinancier{
		repo: repo,
	}
}

func (s *dobbyFinancier) CreatePeriod(ctx context.Context, start, end time.Time) (*Period, error) {
	p := &Period{
		ID:        uuid.New(),
		StartDate: start,
		EndDate:   end,
	}
	if err := s.repo.SavePeriod(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *dobbyFinancier) GetPeriodSummary(ctx context.Context, id uuid.UUID) (*PeriodSummary, error) {
	period, err := s.repo.GetPeriod(ctx, id)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.ListTransactions(ctx, TransactionFilter{PeriodID: &id})
	if err != nil {
		return nil, err
	}

	envelopes, err := s.repo.ListEnvelopes(ctx)
	if err != nil {
		return nil, err
	}

	summary := &PeriodSummary{
		Period: *period,
	}

	envelopeStatsMap := make(map[uuid.UUID]*EnvelopeStat)
	for _, env := range envelopes {
		envelopeStatsMap[env.ID] = &EnvelopeStat{
			Envelope: env,
		}
	}

	for _, t := range transactions {
		stat, ok := envelopeStatsMap[t.EnvelopeID]
		if !ok {
			continue
		}

		if t.Amount > 0 {
			stat.Allocated += t.Amount
			summary.TotalBudget += t.Amount
		} else {
			stat.Spent += t.Amount // Amount is negative for expenses
			summary.TotalSpent += t.Amount
		}
	}

	for _, env := range envelopes {
		stat := envelopeStatsMap[env.ID]
		stat.Remaining = stat.Allocated + stat.Spent
		summary.EnvelopeStats = append(summary.EnvelopeStats, *stat)
	}

	summary.TotalRemaining = summary.TotalBudget + summary.TotalSpent

	return summary, nil
}

func (s *dobbyFinancier) ListPeriods(ctx context.Context) ([]Period, error) {
	return s.repo.ListPeriods(ctx)
}

func (s *dobbyFinancier) RecordTransaction(ctx context.Context, t Transaction) (*Transaction, error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if err := s.repo.SaveTransaction(ctx, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *dobbyFinancier) CreateEnvelope(ctx context.Context, name string) (*Envelope, error) {
	e := &Envelope{
		ID:   uuid.New(),
		Name: name,
	}
	if err := s.repo.SaveEnvelope(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *dobbyFinancier) ListEnvelopes(ctx context.Context) ([]Envelope, error) {
	return s.repo.ListEnvelopes(ctx)
}

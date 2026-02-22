package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type dobbyFinancier struct {
	repo      Repository
	txManager TransactionManager
}

func NewDobbyFinancier(repo Repository, txManager TransactionManager) FinanceService {
	return &dobbyFinancier{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *dobbyFinancier) CreatePeriod(ctx context.Context, start, end *time.Time) (*Period, error) {
	if start == nil {
		startTime, err := calculateNextPeriodStartTime(time.Now().AddDate(0, -1, 0))
		if err != nil {
			return nil, err
		}
		start = &startTime
	}
	if end == nil {
		endTime, err := calculateNextPeriodStartTime(time.Now())
		if err != nil {
			return nil, err
		}
		end = &endTime
	}
	p := &Period{
		ID:        uuid.New(),
		StartDate: *start,
		EndDate:   *end,
	}
	if err := s.repo.SavePeriod(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func calculateNextPeriodStartTime(referenceTime time.Time) (time.Time, error) {
	y, m, _ := referenceTime.Date()
	nextPeriodStart := time.Date(y, time.Month(m+1), 5, 0, 0, 0, 0, &time.Location{})
	weekday := nextPeriodStart.Weekday()
	moveByDays := 0
	if weekday > 5 {
		moveByDays = int(weekday) - 5
	}
	if moveByDays == 0 {
		return nextPeriodStart, nil
	}
	return nextPeriodStart.AddDate(0, 0, -moveByDays), nil
}

func (s *dobbyFinancier) GetCurrentPeriod(ctx context.Context) (*PeriodSummary, error) {
	p, err := s.repo.GetCurrentPeriod(ctx)
	if errors.Is(err, ErrNotFound) {
		slog.Warn("Failed to get current period, creating a new one", "error", err)
		p, err = s.CreatePeriod(ctx, nil, nil)
		if err != nil {
			return nil, err
		}
		return s.GetPeriodSummary(ctx, p.ID)
	} else if err != nil {
		return nil, err
	}
	return s.GetPeriodSummary(ctx, p.ID)
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

	err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		return s.repo.SaveTransaction(ctx, &t)
	})

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *dobbyFinancier) ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	return s.repo.ListTransactions(ctx, filter)
}

func (s *dobbyFinancier) GetTransaction(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *dobbyFinancier) UpdateTransaction(ctx context.Context, t Transaction) (*Transaction, error) {
	_, err := s.repo.GetTransaction(ctx, t.ID)
	if err != nil {
		return nil, err
	}

	err = s.txManager.WithTx(ctx, func(ctx context.Context) error {
		return s.repo.SaveTransaction(ctx, &t)
	})
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *dobbyFinancier) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTransaction(ctx, id)
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

func (s *dobbyFinancier) DeleteEnvelope(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEnvelope(ctx, id)
}

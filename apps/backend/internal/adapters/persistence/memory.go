package persistence

import (
	"context"
	"sync"

	"github.com/ChaPerx64/dobby/apps/backend/internal/service"
	"github.com/google/uuid"
)

type memoryRepo struct {
	mu           sync.RWMutex
	users        map[uuid.UUID]service.User
	periods      map[uuid.UUID]service.Period
	envelopes    map[uuid.UUID]service.Envelope
	transactions map[uuid.UUID]service.Transaction
}

func NewMemoryRepository() service.Repository {
	return &memoryRepo{
		users:        make(map[uuid.UUID]service.User),
		periods:      make(map[uuid.UUID]service.Period),
		envelopes:    make(map[uuid.UUID]service.Envelope),
		transactions: make(map[uuid.UUID]service.Transaction),
	}
}

func (r *memoryRepo) WithTx(ctx context.Context, fn func(repo service.Repository) error) error {
	// Memory repo doesn't support real transactions, just call the function
	return fn(r)
}

func (r *memoryRepo) SaveUser(ctx context.Context, u *service.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = *u
	return nil
}

func (r *memoryRepo) GetUser(ctx context.Context, id uuid.UUID) (*service.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, service.ErrNotFound
	}
	return &u, nil
}

func (r *memoryRepo) SavePeriod(ctx context.Context, p *service.Period) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.periods[p.ID] = *p
	return nil
}

func (r *memoryRepo) GetPeriod(ctx context.Context, id uuid.UUID) (*service.Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.periods[id]
	if !ok {
		return nil, service.ErrNotFound
	}
	return &p, nil
}

func (r *memoryRepo) ListPeriods(ctx context.Context) ([]service.Period, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]service.Period, 0, len(r.periods))
	for _, p := range r.periods {
		res = append(res, p)
	}
	return res, nil
}

func (r *memoryRepo) SaveEnvelope(ctx context.Context, e *service.Envelope) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.envelopes[e.ID] = *e
	return nil
}

func (r *memoryRepo) ListEnvelopes(ctx context.Context) ([]service.Envelope, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]service.Envelope, 0, len(r.envelopes))
	for _, e := range r.envelopes {
		res = append(res, e)
	}
	return res, nil
}

func (r *memoryRepo) SaveTransaction(ctx context.Context, t *service.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transactions[t.ID] = *t
	return nil
}

func (r *memoryRepo) ListTransactions(ctx context.Context, filter service.TransactionFilter) ([]service.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]service.Transaction, 0)
	for _, t := range r.transactions {
		if filter.PeriodID != nil && t.PeriodID != *filter.PeriodID {
			continue
		}
		if filter.EnvelopeID != nil && t.EnvelopeID != *filter.EnvelopeID {
			continue
		}
		res = append(res, t)
	}
	return res, nil
}

package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ChaPerx64/dobby/apps/backend/internal/service"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type uowKey struct{}

type psqlRepo struct {
	db *sql.DB
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func NewPostgresRepository(db *sql.DB) service.Repository {
	return &psqlRepo{db: db}
}

func (r *psqlRepo) getDB(ctx context.Context) DB {
	if tx, ok := ctx.Value(uowKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.db
}

type psqlTxManager struct {
	db *sql.DB
}

func NewPostgresTransactionManager(db *sql.DB) service.TransactionManager {
	return &psqlTxManager{db: db}
}

func (m *psqlTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	ctxWithTx := context.WithValue(ctx, uowKey{}, tx)
	if err := fn(ctxWithTx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *psqlRepo) SaveUser(ctx context.Context, u *service.User) error {
	query := `INSERT INTO users (id, name) VALUES ($1, $2)
              ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`
	_, err := r.getDB(ctx).ExecContext(ctx, query, u.ID, u.Name)
	return err
}

func (r *psqlRepo) GetUser(ctx context.Context, id uuid.UUID) (*service.User, error) {
	query := `SELECT id, name FROM users WHERE id = $1`
	u := &service.User{}
	err := r.getDB(ctx).QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name)
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return u, err
}

func (r *psqlRepo) SavePeriod(ctx context.Context, p *service.Period) error {
	query := `INSERT INTO financial_periods (id, start_dt, end_dt) VALUES ($1, $2, $3)
              ON CONFLICT (id) DO UPDATE SET start_dt = EXCLUDED.start_dt, end_dt = EXCLUDED.end_dt`
	_, err := r.getDB(ctx).ExecContext(ctx, query, p.ID, p.StartDate, p.EndDate)
	return err
}

func (r *psqlRepo) GetPeriod(ctx context.Context, id uuid.UUID) (*service.Period, error) {
	query := `SELECT id, start_dt, end_dt FROM financial_periods WHERE id = $1`
	p := &service.Period{}
	err := r.getDB(ctx).QueryRowContext(ctx, query, id).Scan(&p.ID, &p.StartDate, &p.EndDate)
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return p, err
}

func (r *psqlRepo) GetCurrentPeriod(ctx context.Context) (*service.Period, error) {
	query := `SELECT id, start_dt, end_dt FROM financial_periods 
              WHERE NOW() BETWEEN start_dt AND end_dt 
              ORDER BY start_dt ASC LIMIT 1`
	p := &service.Period{}
	err := r.getDB(ctx).QueryRowContext(ctx, query).Scan(&p.ID, &p.StartDate, &p.EndDate)
	if err == sql.ErrNoRows {
		return nil, service.ErrNotFound
	}
	return p, err
}

func (r *psqlRepo) ListPeriods(ctx context.Context) ([]service.Period, error) {
	query := `SELECT id, start_dt, end_dt FROM financial_periods ORDER BY start_dt DESC`
	rows, err := r.getDB(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []service.Period
	for rows.Next() {
		var p service.Period
		if err := rows.Scan(&p.ID, &p.StartDate, &p.EndDate); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (r *psqlRepo) SaveEnvelope(ctx context.Context, e *service.Envelope) error {
	query := `INSERT INTO envelopes (id, user_id, name) VALUES ($1, $2, $3)
              ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, user_id = EXCLUDED.user_id`
	_, err := r.getDB(ctx).ExecContext(ctx, query, e.ID, e.UserID, e.Name)
	return err
}

func (r *psqlRepo) ListEnvelopes(ctx context.Context) ([]service.Envelope, error) {
	query := `SELECT id, user_id, name FROM envelopes`
	rows, err := r.getDB(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []service.Envelope
	for rows.Next() {
		var e service.Envelope
		if err := rows.Scan(&e.ID, &e.UserID, &e.Name); err != nil {
			return nil, err
		}
		res = append(res, e)
	}
	return res, nil
}

func (r *psqlRepo) SaveTransaction(ctx context.Context, t *service.Transaction) error {
	query := `INSERT INTO transactions (id, financial_period_id, user_id, envelope_id, category, amount, description, date)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
              ON CONFLICT (id) DO UPDATE SET 
                financial_period_id = EXCLUDED.financial_period_id,
                user_id = EXCLUDED.user_id,
                envelope_id = EXCLUDED.envelope_id,
                category = EXCLUDED.category,
                amount = EXCLUDED.amount,
                description = EXCLUDED.description,
                date = EXCLUDED.date`
	_, err := r.getDB(ctx).ExecContext(ctx, query, t.ID, t.PeriodID, t.UserID, t.EnvelopeID, t.Category, t.Amount, t.Description, t.Date)
	return err
}

func (r *psqlRepo) ListTransactions(ctx context.Context, filter service.TransactionFilter) ([]service.Transaction, error) {
	query := `SELECT id, financial_period_id, user_id, envelope_id, category, amount, description, date FROM transactions WHERE 1=1`
	var args []interface{}
	argCount := 1

	if filter.PeriodID != nil {
		query += fmt.Sprintf(" AND financial_period_id = $%d", argCount)
		args = append(args, *filter.PeriodID)
		argCount++
	}
	if filter.EnvelopeID != nil {
		query += fmt.Sprintf(" AND envelope_id = $%d", argCount)
		args = append(args, *filter.EnvelopeID)
		argCount++
	}

	query += " ORDER BY date DESC"

	rows, err := r.getDB(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []service.Transaction
	for rows.Next() {
		var t service.Transaction
		if err := rows.Scan(&t.ID, &t.PeriodID, &t.UserID, &t.EnvelopeID, &t.Category, &t.Amount, &t.Description, &t.Date); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

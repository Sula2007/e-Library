package postgres

import (
	"context"
	"database/sql"

	"github.com/Sula2007/payment-service/internal/domain"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	query := `INSERT INTO payments (id, user_id, book_id, amount, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.UserID, p.BookID, p.Amount, p.Status, p.CreatedAt)
	return err
}

func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	query := `SELECT id, user_id, book_id, amount, status, created_at FROM payments WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	p := &domain.Payment{}
	err := row.Scan(&p.ID, &p.UserID, &p.BookID, &p.Amount, &p.Status, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PaymentRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Payment, error) {
	query := `SELECT id, user_id, book_id, amount, status, created_at FROM payments WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []*domain.Payment
	for rows.Next() {
		p := &domain.Payment{}
		if err := rows.Scan(&p.ID, &p.UserID, &p.BookID, &p.Amount, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id, status string) error {
	query := `UPDATE payments SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
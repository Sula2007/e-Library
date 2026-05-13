package usecase

import (
	"context"
	"time"

	"github.com/Sula2007/payment-service/internal/domain"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) error
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Payment, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type PaymentUsecase struct {
	repo PaymentRepository
}

func NewPaymentUsecase(repo PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: repo}
}

func (u *PaymentUsecase) CreatePayment(ctx context.Context, userID, bookID string, amount float64) (*domain.Payment, error) {
	p := &domain.Payment{
		ID:        uuid.New().String(),
		UserID:    userID,
		BookID:    bookID,
		Amount:    amount,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}
	if err := u.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (u *PaymentUsecase) GetPayment(ctx context.Context, id string) (*domain.Payment, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PaymentUsecase) GetUserPayments(ctx context.Context, userID string) ([]*domain.Payment, error) {
	return u.repo.GetByUserID(ctx, userID)
}

func (u *PaymentUsecase) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	return u.repo.UpdateStatus(ctx, id, status)
}
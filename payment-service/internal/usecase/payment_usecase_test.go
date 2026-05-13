package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Sula2007/payment-service/internal/domain"
	"github.com/Sula2007/payment-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) GetByUserID(ctx context.Context, userID string) ([]*domain.Payment, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepo) UpdateStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestCreatePayment_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	uc := usecase.NewPaymentUsecase(repo)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(nil)
	p, err := uc.CreatePayment(context.Background(), "user-1", "book-1", 9.99)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "user-1", p.UserID)
	assert.Equal(t, "book-1", p.BookID)
	assert.Equal(t, 9.99, p.Amount)
	assert.Equal(t, domain.StatusPending, p.Status)
	repo.AssertExpectations(t)
}

func TestCreatePayment_RepoError(t *testing.T) {
	repo := new(MockPaymentRepo)
	uc := usecase.NewPaymentUsecase(repo)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Payment")).Return(errors.New("db error"))
	p, err := uc.CreatePayment(context.Background(), "user-1", "book-1", 9.99)
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestGetPayment_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	uc := usecase.NewPaymentUsecase(repo)
	expected := &domain.Payment{ID: "pay-1", UserID: "user-1", Amount: 9.99, Status: domain.StatusPending}
	repo.On("GetByID", mock.Anything, "pay-1").Return(expected, nil)
	p, err := uc.GetPayment(context.Background(), "pay-1")
	assert.NoError(t, err)
	assert.Equal(t, "pay-1", p.ID)
}

func TestGetPayment_NotFound(t *testing.T) {
	repo := new(MockPaymentRepo)
	uc := usecase.NewPaymentUsecase(repo)
	repo.On("GetByID", mock.Anything, "not-exist").Return(nil, errors.New("not found"))
	p, err := uc.GetPayment(context.Background(), "not-exist")
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestGetUserPayments_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	uc := usecase.NewPaymentUsecase(repo)
	payments := []*domain.Payment{
		{ID: "pay-1", UserID: "user-1", Amount: 9.99},
		{ID: "pay-2", UserID: "user-1", Amount: 4.99},
	}
	repo.On("GetByUserID", mock.Anything, "user-1").Return(payments, nil)
	result, err := uc.GetUserPayments(context.Background(), "user-1")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUpdatePaymentStatus_Success(t *testing.T) {
	repo := new(MockPaymentRepo)
	uc := usecase.NewPaymentUsecase(repo)
	repo.On("UpdateStatus", mock.Anything, "pay-1", domain.StatusPaid).Return(nil)
	err := uc.UpdatePaymentStatus(context.Background(), "pay-1", domain.StatusPaid)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Sula2007/user-service/internal/domain"
	"github.com/Sula2007/user-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserCache struct {
	mock.Mock
}

func (m *MockUserCache) Set(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserCache) Get(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserCache) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	repo := new(MockUserRepo)
	cache := new(MockUserCache)
	uc := usecase.NewUserUsecase(repo, cache)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	user, err := uc.Register(context.Background(), "Sultan", "sultan@gmail.com", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Sultan", user.Name)
	assert.Equal(t, "sultan@gmail.com", user.Email)
	repo.AssertExpectations(t)
}

func TestRegister_RepoError(t *testing.T) {
	repo := new(MockUserRepo)
	cache := new(MockUserCache)
	uc := usecase.NewUserUsecase(repo, cache)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("db error"))
	user, err := uc.Register(context.Background(), "Sultan", "sultan@gmail.com", "password123")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(MockUserRepo)
	cache := new(MockUserCache)
	uc := usecase.NewUserUsecase(repo, cache)
	repo.On("GetByEmail", mock.Anything, "notfound@gmail.com").Return(nil, errors.New("not found"))
	user, err := uc.Login(context.Background(), "notfound@gmail.com", "password")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetProfile_FromCache(t *testing.T) {
	repo := new(MockUserRepo)
	cache := new(MockUserCache)
	uc := usecase.NewUserUsecase(repo, cache)
	cachedUser := &domain.User{ID: "some-id", Name: "Sultan", Email: "sultan@gmail.com"}
	cache.On("Get", mock.Anything, "some-id").Return(cachedUser, nil)
	user, err := uc.GetProfile(context.Background(), "some-id")
	assert.NoError(t, err)
	assert.Equal(t, "Sultan", user.Name)
	repo.AssertNotCalled(t, "GetByID")
}

func TestDeleteUser_Success(t *testing.T) {
	repo := new(MockUserRepo)
	cache := new(MockUserCache)
	uc := usecase.NewUserUsecase(repo, cache)
	repo.On("Delete", mock.Anything, "some-id").Return(nil)
	cache.On("Delete", mock.Anything, "some-id").Return(nil)
	err := uc.DeleteUser(context.Background(), "some-id")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
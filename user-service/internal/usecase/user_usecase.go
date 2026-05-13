package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Sula2007/user-service/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type UserCache interface {
	Set(ctx context.Context, user *domain.User) error
	Get(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type UserUsecase struct {
	repo  UserRepository
	cache UserCache
}

func NewUserUsecase(repo UserRepository, cache UserCache) *UserUsecase {
	return &UserUsecase{repo: repo, cache: cache}
}

func (u *UserUsecase) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func (u *UserUsecase) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	cached, err := u.cache.Get(ctx, id)
	if err == nil {
		return cached, nil
	}
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u.cache.Set(ctx, user)
	return user, nil
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, id, name, email string) error {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.Name = name
	user.Email = email
	if err := u.repo.Update(ctx, user); err != nil {
		return err
	}
	u.cache.Delete(ctx, id)
	return nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}
	u.cache.Delete(ctx, id)
	return nil
}
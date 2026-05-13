package grpc

import (
	"context"

	"github.com/Sula2007/user-service/internal/usecase"
	gen "github.com/Sula2007/user-service/proto/gen"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type UserHandler struct {
	gen.UnimplementedUserServiceServer
	usecase   *usecase.UserUsecase
	jwtSecret string
}

func NewUserHandler(usecase *usecase.UserUsecase, jwtSecret string) *UserHandler {
	return &UserHandler{usecase: usecase, jwtSecret: jwtSecret}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *gen.RegisterUserRequest) (*gen.RegisterUserResponse, error) {
	user, err := h.usecase.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.RegisterUserResponse{Id: user.ID, Message: "registered successfully"}, nil
}

func (h *UserHandler) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	user, err := h.usecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.LoginUserResponse{Token: signed, Message: "login successful"}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *gen.GetUserProfileRequest) (*gen.GetUserProfileResponse, error) {
	user, err := h.usecase.GetProfile(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &gen.GetUserProfileResponse{Id: user.ID, Name: user.Name, Email: user.Email}, nil
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *gen.UpdateUserProfileRequest) (*gen.UpdateUserProfileResponse, error) {
	err := h.usecase.UpdateProfile(ctx, req.Id, req.Name, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.UpdateUserProfileResponse{Message: "updated successfully"}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *gen.DeleteUserRequest) (*gen.DeleteUserResponse, error) {
	err := h.usecase.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.DeleteUserResponse{Message: "deleted successfully"}, nil
}
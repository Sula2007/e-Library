package grpc

import (
	"context"

	"github.com/Sula2007/payment-service/internal/usecase"
	gen "github.com/Sula2007/payment-service/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	gen.UnimplementedPaymentServiceServer
	usecase *usecase.PaymentUsecase
}

func NewPaymentHandler(uc *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{usecase: uc}
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *gen.CreatePaymentRequest) (*gen.CreatePaymentResponse, error) {
	p, err := h.usecase.CreatePayment(ctx, req.UserId, req.BookId, req.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.CreatePaymentResponse{Id: p.ID, Status: p.Status, Message: "payment created"}, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *gen.GetPaymentRequest) (*gen.GetPaymentResponse, error) {
	p, err := h.usecase.GetPayment(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &gen.GetPaymentResponse{Payment: &gen.Payment{
		Id:        p.ID,
		UserId:    p.UserID,
		BookId:    p.BookID,
		Amount:    p.Amount,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.String(),
	}}, nil
}

func (h *PaymentHandler) GetUserPayments(ctx context.Context, req *gen.GetUserPaymentsRequest) (*gen.GetUserPaymentsResponse, error) {
	payments, err := h.usecase.GetUserPayments(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var result []*gen.Payment
	for _, p := range payments {
		result = append(result, &gen.Payment{
			Id:        p.ID,
			UserId:    p.UserID,
			BookId:    p.BookID,
			Amount:    p.Amount,
			Status:    p.Status,
			CreatedAt: p.CreatedAt.String(),
		})
	}
	return &gen.GetUserPaymentsResponse{Payments: result}, nil
}

func (h *PaymentHandler) UpdatePaymentStatus(ctx context.Context, req *gen.UpdatePaymentStatusRequest) (*gen.UpdatePaymentStatusResponse, error) {
	err := h.usecase.UpdatePaymentStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &gen.UpdatePaymentStatusResponse{Message: "status updated"}, nil
}
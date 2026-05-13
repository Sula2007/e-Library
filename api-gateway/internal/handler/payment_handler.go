package handler

import (
	"context"
	"net/http"

	gen "github.com/Sula2007/api-gateway/proto/gen"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentHandler struct {
	client gen.PaymentServiceClient
}

func NewPaymentHandler(addr string) (*PaymentHandler, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &PaymentHandler{client: gen.NewPaymentServiceClient(conn)}, nil
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req gen.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.client.CreatePayment(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.GetPayment(context.Background(), &gen.GetPaymentRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) GetUserPayments(c *gin.Context) {
	userID := c.Param("user_id")
	resp, err := h.client.GetUserPayments(context.Background(), &gen.GetUserPaymentsRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	var req gen.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = c.Param("id")
	resp, err := h.client.UpdatePaymentStatus(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
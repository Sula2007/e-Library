package main

import (
	"log"
	"os"

	"github.com/Sula2007/api-gateway/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50051"
	}

	paymentServiceAddr := os.Getenv("PAYMENT_SERVICE_ADDR")
	if paymentServiceAddr == "" {
		paymentServiceAddr = "localhost:50054"
	}

	userHandler, err := handler.NewUserHandler(userServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	paymentHandler, err := handler.NewPaymentHandler(paymentServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}

	r := gin.Default()

	r.POST("/users/register", userHandler.Register)
	r.POST("/users/login", userHandler.Login)
	r.GET("/users/:id", userHandler.GetProfile)
	r.PUT("/users/:id", userHandler.UpdateProfile)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	r.POST("/payments", paymentHandler.CreatePayment)
	r.GET("/payments/:id", paymentHandler.GetPayment)
	r.GET("/payments/user/:user_id", paymentHandler.GetUserPayments)
	r.PUT("/payments/:id/status", paymentHandler.UpdatePaymentStatus)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("API Gateway started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
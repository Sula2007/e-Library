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

	userHandler, err := handler.NewUserHandler(userServiceAddr)
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}

	r := gin.Default()

	r.POST("/users/register", userHandler.Register)
	r.POST("/users/login", userHandler.Login)
	r.GET("/users/:id", userHandler.GetProfile)
	r.PUT("/users/:id", userHandler.UpdateProfile)
	r.DELETE("/users/:id", userHandler.DeleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway started on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
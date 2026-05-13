package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/Sula2007/payment-service/internal/config"
	grpchandler "github.com/Sula2007/payment-service/internal/delivery/grpc"
	postgrepo "github.com/Sula2007/payment-service/internal/repository/postgres"
	"github.com/Sula2007/payment-service/internal/usecase"
	gen "github.com/Sula2007/payment-service/proto/gen"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	m, err := migrate.New("file://migrations", cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to init migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := postgrepo.NewPaymentRepository(db)
	uc := usecase.NewPaymentUsecase(repo)
	handler := grpchandler.NewPaymentHandler(uc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	gen.RegisterPaymentServiceServer(server, handler)

	log.Printf("gRPC server started on port %s", cfg.GRPCPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
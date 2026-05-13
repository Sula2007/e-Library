package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	"github.com/Sula2007/user-service/internal/config"
	grpchandler "github.com/Sula2007/user-service/internal/delivery/grpc"
	"github.com/Sula2007/user-service/internal/email"
	natssub "github.com/Sula2007/user-service/internal/nats"
	postgrepo "github.com/Sula2007/user-service/internal/repository/postgres"
	rediscache "github.com/Sula2007/user-service/internal/repository/redis"
	"github.com/Sula2007/user-service/internal/usecase"
	gen "github.com/Sula2007/user-service/proto/gen"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}
	defer nc.Close()

	repo := postgrepo.NewUserRepository(db)
	cache := rediscache.NewUserCache(redisClient)
	uc := usecase.NewUserUsecase(repo, cache)
	handler := grpchandler.NewUserHandler(uc, cfg.JWTSecret)

	emailSender := email.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass)
	subscriber := natssub.NewSubscriber(nc, repo, emailSender)
	subscriber.Subscribe()

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	gen.RegisterUserServiceServer(server, handler)

	log.Printf("gRPC server started on port %s", cfg.GRPCPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
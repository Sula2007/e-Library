package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Sula2007/payment-service/internal/domain"
	"github.com/Sula2007/payment-service/internal/repository/postgres"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	pgmodule "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	pgContainer, err := pgmodule.Run(ctx,
		"postgres:15",
		pgmodule.WithDatabase("testdb"),
		pgmodule.WithUsername("postgres"),
		pgmodule.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS payments (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		book_id VARCHAR(36) NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db, func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}
}

func TestPaymentRepo_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewPaymentRepository(db)
	p := &domain.Payment{
		ID:        "pay-1",
		UserID:    "user-1",
		BookID:    "book-1",
		Amount:    9.99,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	err := repo.Create(context.Background(), p)
	assert.NoError(t, err)
}

func TestPaymentRepo_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewPaymentRepository(db)
	p := &domain.Payment{
		ID:        "pay-2",
		UserID:    "user-1",
		BookID:    "book-1",
		Amount:    9.99,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	err := repo.Create(context.Background(), p)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), "pay-2")
	assert.NoError(t, err)
	assert.Equal(t, "user-1", found.UserID)
	assert.Equal(t, 9.99, found.Amount)
}

func TestPaymentRepo_GetByUserID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewPaymentRepository(db)
	payments := []*domain.Payment{
		{ID: "pay-3", UserID: "user-2", BookID: "book-1", Amount: 9.99, Status: domain.StatusPending, CreatedAt: time.Now()},
		{ID: "pay-4", UserID: "user-2", BookID: "book-2", Amount: 4.99, Status: domain.StatusPending, CreatedAt: time.Now()},
	}

	for _, p := range payments {
		err := repo.Create(context.Background(), p)
		assert.NoError(t, err)
	}

	result, err := repo.GetByUserID(context.Background(), "user-2")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPaymentRepo_UpdateStatus(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewPaymentRepository(db)
	p := &domain.Payment{
		ID:        "pay-5",
		UserID:    "user-1",
		BookID:    "book-1",
		Amount:    9.99,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	err := repo.Create(context.Background(), p)
	assert.NoError(t, err)

	err = repo.UpdateStatus(context.Background(), "pay-5", domain.StatusPaid)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), "pay-5")
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusPaid, found.Status)
}
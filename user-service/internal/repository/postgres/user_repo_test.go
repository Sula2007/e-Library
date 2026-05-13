package postgres_test

import (
	"context"
	"testing"

	"github.com/Sula2007/user-service/internal/domain"
	"github.com/Sula2007/user-service/internal/repository/postgres"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	pgmodule "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gosql "database/sql"
	"time"
)

func setupTestDB(t *testing.T) (*gosql.DB, func()) {
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

	db, err := gosql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
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

func TestUserRepo_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewUserRepository(db)
	user := &domain.User{
		ID:           "test-id-1",
		Name:         "Sultan",
		Email:        "sultan@test.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
}

func TestUserRepo_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewUserRepository(db)
	user := &domain.User{
		ID:           "test-id-2",
		Name:         "Sultan",
		Email:        "sultan2@test.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), "test-id-2")
	assert.NoError(t, err)
	assert.Equal(t, "Sultan", found.Name)
	assert.Equal(t, "sultan2@test.com", found.Email)
}

func TestUserRepo_GetByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewUserRepository(db)
	user := &domain.User{
		ID:           "test-id-3",
		Name:         "Sultan",
		Email:        "sultan3@test.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	found, err := repo.GetByEmail(context.Background(), "sultan3@test.com")
	assert.NoError(t, err)
	assert.Equal(t, "test-id-3", found.ID)
}

func TestUserRepo_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewUserRepository(db)
	user := &domain.User{
		ID:           "test-id-4",
		Name:         "Sultan",
		Email:        "sultan4@test.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	user.Name = "Sultan Updated"
	err = repo.Update(context.Background(), user)
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), "test-id-4")
	assert.NoError(t, err)
	assert.Equal(t, "Sultan Updated", found.Name)
}

func TestUserRepo_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewUserRepository(db)
	user := &domain.User{
		ID:           "test-id-5",
		Name:         "Sultan",
		Email:        "sultan5@test.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), "test-id-5")
	assert.NoError(t, err)

	found, err := repo.GetByID(context.Background(), "test-id-5")
	assert.Error(t, err)
	assert.Nil(t, found)
}
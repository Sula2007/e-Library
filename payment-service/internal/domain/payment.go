package domain

import "time"

type Payment struct {
	ID        string
	UserID    string
	BookID    string
	Amount    float64
	Status    string
	CreatedAt time.Time
}

const (
	StatusPending = "pending"
	StatusPaid    = "paid"
	StatusFailed  = "failed"
)
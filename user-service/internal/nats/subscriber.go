package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Sula2007/user-service/internal/email"
	"github.com/Sula2007/user-service/internal/repository/postgres"
	"github.com/nats-io/nats.go"
)

type BorrowEvent struct {
	UserID    string `json:"user_id"`
	BookTitle string `json:"book_title"`
}

type Subscriber struct {
	nc     *nats.Conn
	repo   *postgres.UserRepository
	sender *email.Sender
}

func NewSubscriber(nc *nats.Conn, repo *postgres.UserRepository, sender *email.Sender) *Subscriber {
	return &Subscriber{nc: nc, repo: repo, sender: sender}
}

func (s *Subscriber) Subscribe() {
	s.nc.Subscribe("borrow.created", func(msg *nats.Msg) {
		var event BorrowEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("failed to unmarshal borrow.created: %v", err)
			return
		}
		user, err := s.repo.GetByID(context.Background(), event.UserID)
		if err != nil {
			log.Printf("user not found: %v", err)
			return
		}
		if err := s.sender.Send(user.Email, "Book Borrowed", "You have borrowed the book: "+event.BookTitle); err != nil {
			log.Printf("failed to send email: %v", err)
		}
	})

	s.nc.Subscribe("borrow.returned", func(msg *nats.Msg) {
		var event BorrowEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("failed to unmarshal borrow.returned: %v", err)
			return
		}
		user, err := s.repo.GetByID(context.Background(), event.UserID)
		if err != nil {
			log.Printf("user not found: %v", err)
			return
		}
		if err := s.sender.Send(user.Email, "Book Returned", "You have returned the book: "+event.BookTitle); err != nil {
			log.Printf("failed to send email: %v", err)
		}
	})

	log.Println("NATS subscriber started")
}
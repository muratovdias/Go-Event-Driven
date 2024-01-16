package booking

import (
	"context"
	"github.com/google/uuid"
	"tickets/internal/entities"
	"tickets/internal/repository"
)

type Service struct {
	repo repository.Booking
}

func NewService(repo repository.Booking) *Service {
	return &Service{repo: repo}
}

func (s *Service) BookTicket(ctx context.Context, booking entities.Booking) (string, error) {
	booking.BookingID = uuid.New()
	return s.repo.BookTicket(ctx, booking)
}

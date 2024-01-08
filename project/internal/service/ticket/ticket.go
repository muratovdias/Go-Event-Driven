package ticket

import (
	"context"
	"tickets/internal/entities"
	"tickets/internal/repository"
)

type Service struct {
	repo repository.Ticket
}

func NewService(repo repository.Ticket) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveTicket(ctx context.Context, ticket entities.TicketBookingConfirmed) error {
	return s.repo.SaveTicket(ctx, ticket)
}

func (s *Service) DeleteTicket(ctx context.Context, ticketID string) error {
	return s.repo.DeleteTicket(ctx, ticketID)
}

func (s *Service) TicketList(ctx context.Context) ([]entities.TicketList, error) {
	return s.repo.TicketList(ctx)
}

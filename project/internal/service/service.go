package service

import (
	"context"
	"tickets/internal/entities"
	"tickets/internal/repository"
	"tickets/internal/service/booking"
	"tickets/internal/service/show"
	"tickets/internal/service/ticket"
)

type ReceiptsClient interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type SpreadsheetsClient interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}

type FilesClient interface {
	StoreTicketContent(ctx context.Context, ticket entities.TicketBookingConfirmed) error
}

type Ticket interface {
	SaveTicket(ctx context.Context, ticket entities.TicketBookingConfirmed) error
	DeleteTicket(ctx context.Context, ticketID string) error
	TicketList(ctx context.Context) ([]entities.TicketList, error)
}

type Show interface {
	NewShow(ctx context.Context, show entities.Show) (string, error)
}

type Booking interface {
	BookTicket(ctx context.Context, booking entities.Booking) (string, error)
}

type Service struct {
	ReceiptsClient
	SpreadsheetsClient
	FilesClient
	Ticket
	Show
	Booking
}

func NewService(receiptsClient ReceiptsClient, spreadsheetsClient SpreadsheetsClient, filesClient FilesClient,
	repo *repository.Repository) *Service {

	return &Service{
		ReceiptsClient:     receiptsClient,
		SpreadsheetsClient: spreadsheetsClient,
		FilesClient:        filesClient,
		Ticket:             ticket.NewService(repo.Ticket),
		Show:               show.NewService(repo.Show),
		Booking:            booking.NewService(repo.Booking),
	}

}

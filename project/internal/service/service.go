package service

import (
	"context"
	"github.com/google/uuid"
	"tickets/internal/entities"
	"tickets/internal/repository"
	"tickets/internal/service/booking"
	"tickets/internal/service/show"
	"tickets/internal/service/ticket"
)

type ReceiptsClient interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
	// TODO: void the receipt
}

type SpreadsheetsClient interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}

type FilesClient interface {
	StoreTicketContent(ctx context.Context, ticket entities.TicketBookingConfirmed) error
}

type DeadNationClient interface {
	BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error
}

type Ticket interface {
	SaveTicket(ctx context.Context, ticket entities.TicketBookingConfirmed) error
	DeleteTicket(ctx context.Context, ticketID string) error
	TicketList(ctx context.Context) ([]entities.TicketList, error)
}

type Show interface {
	NewShow(ctx context.Context, show entities.Show) (string, error)
	ShowByID(ctx context.Context, showId uuid.UUID) (entities.Show, error)
}

type Booking interface {
	BookTicket(ctx context.Context, booking entities.Booking) (string, error)
}

type Service struct {
	ReceiptsClient
	SpreadsheetsClient
	FilesClient
	DeadNationClient
	Ticket
	Show
	Booking
	// TODO: payment service
}

func NewService(receiptsClient ReceiptsClient,
	spreadsheetsClient SpreadsheetsClient,
	filesClient FilesClient,
	deadNationClient DeadNationClient,
	repo *repository.Repository) *Service {

	return &Service{
		ReceiptsClient:     receiptsClient,
		SpreadsheetsClient: spreadsheetsClient,
		FilesClient:        filesClient,
		DeadNationClient:   deadNationClient,
		Ticket:             ticket.NewService(repo.Ticket),
		Show:               show.NewService(repo.Show),
		Booking:            booking.NewService(repo.Booking),
	}

}

package broker

import (
	"context"
	"github.com/google/uuid"
	"tickets/internal/entities"
)

type serviceI interface {
	// Receipts
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)

	// Spreadsheets
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error

	// Files
	StoreTicketContent(ctx context.Context, ticket entities.TicketBookingConfirmed) error

	// Dead Nation
	BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error

	// Show
	ShowByID(ctx context.Context, showId uuid.UUID) (entities.Show, error)

	// Ticket
	SaveTicket(ctx context.Context, ticket entities.TicketBookingConfirmed) error
	DeleteTicket(ctx context.Context, ticketID string) error
	TicketList(ctx context.Context) ([]entities.TicketList, error)
}

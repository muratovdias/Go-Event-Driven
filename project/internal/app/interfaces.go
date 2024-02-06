package app

import (
	"context"
	"tickets/internal/entities"
)

type receiptsClient interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
	PutVoidReceiptWithResponse(ctx context.Context, command entities.RefundTicket) error
}

type spreadsheetsClient interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}

type filesClient interface {
	StoreTicketContent(ctx context.Context, ticket entities.TicketBookingConfirmed) error
}

type deadNationClient interface {
	BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error
}

type paymentClient interface {
	PutRefundsWithResponse(ctx context.Context, command entities.RefundTicket) error
}

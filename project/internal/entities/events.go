package entities

import (
	"time"

	"github.com/google/uuid"
)

type EventHeader struct {
	ID          string    `json:"id"`
	PublishedAt time.Time `json:"published_at"`
}

func NewEventHeader() EventHeader {
	return EventHeader{
		ID:          uuid.NewString(),
		PublishedAt: time.Now().UTC(),
	}
}

type TicketBookingConfirmed struct {
	Header EventHeader `json:"header"`

	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`

	BookingID string `json:"booking_id"`
}

func (t *TicketBookingConfirmed) ToSpreadsheetTicketPayload() []string {
	return []string{t.TicketID, t.CustomerEmail, t.Price.Amount, t.Price.Currency}
}

func (t *TicketBookingConfirmed) ToIssueReceiptPayload() IssueReceiptRequest {
	return IssueReceiptRequest{
		TicketID: t.TicketID,
		Price: Money{
			Amount:   t.Price.Amount,
			Currency: t.Price.Currency,
		},
	}
}

type TicketBookingCanceled struct {
	Header EventHeader `json:"header"`

	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

func (t *TicketBookingCanceled) ToSpreadsheetTicketPayload() []string {
	return []string{t.TicketID, t.CustomerEmail, t.Price.Amount, t.Price.Currency}
}

type TicketRefunded struct {
	Header EventHeader `json:"header"`

	TicketID string `json:"ticket_id"`
}

type BookingMade struct {
	Header EventHeader `json:"header"`

	NumberOfTickets int `json:"number_of_tickets"`

	BookingID uuid.UUID `json:"booking_id"`

	CustomerEmail string    `json:"customer_email"`
	ShowId        uuid.UUID `json:"show_id"`
}

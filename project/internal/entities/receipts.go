package entities

import "time"

type VoidReceipt struct {
	TicketID       string
	Reason         string
	IdempotencyKey string
}

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"number"`
	IssuedAt      time.Time `json:"issued_at"`
}

func (t *Ticket) ToIssueReceiptPayload() IssueReceiptRequest {
	return IssueReceiptRequest{
		TicketID: t.TicketID,
		Price: Money{
			Amount:   t.Price.Amount,
			Currency: t.Price.Currency,
		},
	}
}

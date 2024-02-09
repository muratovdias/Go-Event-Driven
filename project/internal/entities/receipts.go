package entities

import (
	"time"
)

type VoidReceipt struct {
	TicketID       string
	Reason         string
	IdempotencyKey string
}

type IssueReceiptRequest struct {
	TicketID       string `json:"ticket_id"`
	IdempotencyKey string `json:"idempotency_key"`
	Price          Money  `json:"price"`
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"number"`
	IssuedAt      time.Time `json:"issued_at"`
}

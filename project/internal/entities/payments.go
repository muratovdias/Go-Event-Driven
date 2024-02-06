package entities

type PaymentRefund struct {
	TicketID       string
	RefundReason   string
	IdempotencyKey string
}

type PaymentRefundRequest struct {
	TicketID       string
	RefundReason   string
	IdempotencyKey string
}

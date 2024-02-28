package command

import (
	"context"
	"fmt"
	"tickets/internal/entities"
)

func (h *Handler) RefundTicket(ctx context.Context, command *entities.RefundTicket) error {
	idempotencyKey := command.Header.IdempotencyKey
	if idempotencyKey == "" {
		return fmt.Errorf("idempotency key is required")
	}

	if err := h.receiptsServiceClient.PutVoidReceiptWithResponse(ctx, entities.VoidReceipt{
		TicketID:       command.TicketID,
		Reason:         "ticket refunded",
		IdempotencyKey: idempotencyKey,
	}); err != nil {
		return fmt.Errorf("failed to void receipt: %w", err)
	}

	if err := h.paymentsServiceClient.PutRefundsWithResponse(ctx, entities.PaymentRefund{
		TicketID:       command.TicketID,
		RefundReason:   "ticket refunded",
		IdempotencyKey: idempotencyKey,
	}); err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	return nil
}

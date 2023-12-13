package broker

import (
	"context"
	"tickets/internal/entities"
)

type eventHandlers struct {
	ticketHandler *ticketHandler
}

func newEventHandlers(service serviceI) *eventHandlers {
	return &eventHandlers{
		ticketHandler: &ticketHandler{service: service},
	}
}

type ticketHandler struct {
	service serviceI
}

func (t *ticketHandler) IssueReceipt(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	if event.Price.Currency == "" {
		event.Price.Currency = "USD"
	}

	_, err := t.service.IssueReceipt(ctx, event.ToIssueReceiptPayload())
	if err != nil {
		return err
	}

	return nil
}

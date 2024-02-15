package broker

import (
	"context"
	"tickets/internal/entities"
)

type commandHandlers struct {
	ticketHandler *ticketHandler
}

func newCommandHandlers(service serviceI, publisher eventPublisher) *commandHandlers {
	return &commandHandlers{
		ticketHandler: &ticketHandler{
			service:        service,
			eventPublisher: publisher,
		},
	}
}

func (t *ticketHandler) RefundTicket(ctx context.Context, command *entities.RefundTicket) error {
	if err := t.service.PutRefundsWithResponse(ctx, *command); err != nil {
		return err
	}

	if err := t.service.PutVoidReceiptWithResponse(ctx, *command); err != nil {
		return err
	}
	return nil
}

package broker

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
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

func (t *ticketHandler) TicketToPrint(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	if event.Price.Currency == "" {
		event.Price.Currency = "USD"
	}

	// add ticket
	if err := t.service.AppendRow(ctx, "tickets-to-print", event.ToSpreadsheetTicketPayload()); err != nil {
		return err
	}

	return nil
}

func (t *ticketHandler) TicketToRefund(ctx context.Context, event *entities.TicketBookingCanceled) error {
	if event.Price.Currency == "" {
		event.Price.Currency = "USD"
	}

	if err := t.service.AppendRow(ctx, "tickets-to-refund", event.ToSpreadsheetTicketPayload()); err != nil {
		return err
	}
	return nil
}

func (t *ticketHandler) ticketHandlers() []cqrs.EventHandler {
	return []cqrs.EventHandler{
		cqrs.NewEventHandler("receipts", t.IssueReceipt),
		cqrs.NewEventHandler("ticket to print", t.TicketToPrint),
		cqrs.NewEventHandler("refund ticket", t.TicketToRefund),
	}
}

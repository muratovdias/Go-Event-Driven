package broker

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
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

func (t *ticketHandler) SaveTicketInDB(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	if err := t.service.SaveTicket(ctx, *event); err != nil {
		return err
	}

	return nil
}

func (t *ticketHandler) StoreTicketContent(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	if err := t.service.StoreTicketContent(ctx, *event); err != nil {
		return err
	}

	return nil
}

func (t *ticketHandler) DeleteTicket(ctx context.Context, event *entities.TicketBookingCanceled) error {
	if err := t.service.DeleteTicket(ctx, event.TicketID); err != nil {
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

func (t *ticketHandler) BookPlaceInDeadNation(ctx context.Context, event *entities.BookingMade) error {
	log.FromContext(ctx).Info("Booking ticket in Dead Nation")

	show, err := t.service.ShowByID(ctx, event.ShowId)
	if err != nil {
		return fmt.Errorf("failed to get show: %w", err)
	}

	err = t.service.BookInDeadNation(ctx, entities.DeadNationBooking{
		CustomerEmail:     event.CustomerEmail,
		DeadNationEventID: show.DeadNationID,
		NumberOfTickets:   event.NumberOfTickets,
		BookingID:         event.BookingID,
	})
	if err != nil {
		return fmt.Errorf("failed to book in dead nation: %w", err)
	}

	return nil
}

func (t *ticketHandler) ticketHandlers() []cqrs.EventHandler {
	return []cqrs.EventHandler{
		cqrs.NewEventHandler("receipts", t.IssueReceipt),
		cqrs.NewEventHandler("ticket to print", t.TicketToPrint),
		cqrs.NewEventHandler("refund ticket", t.TicketToRefund),
		cqrs.NewEventHandler("save ticket in DB", t.SaveTicketInDB),
		cqrs.NewEventHandler("delete ticket from DB", t.DeleteTicket),
		cqrs.NewEventHandler("store ticket content", t.StoreTicketContent),
	}
}

package event

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/sirupsen/logrus"
	"tickets/internal/entities"
)

func (h *Handler) TicketToPrint(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	if event.Price.Currency == "" {
		event.Price.Currency = "USD"
	}

	// add ticket
	if err := h.spreadsheetsService.AppendRow(ctx, "tickets-to-print", event.ToSpreadsheetTicketPayload()); err != nil {
		return err
	}

	return nil
}

func (h *Handler) SaveTicketInDB(ctx context.Context, event *entities.TicketBookingConfirmed) error {

	if err := h.ticketService.SaveTicket(ctx, *event); err != nil {
		return err
	}

	return nil
}

func (h *Handler) StoreTicketContent(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	if err := h.filesAPI.StoreTicketContent(ctx, *event); err != nil {
		return err
	}

	return nil
}

func (h *Handler) DeleteTicket(ctx context.Context, event *entities.TicketBookingCanceled) error {
	if err := h.ticketService.DeleteTicket(ctx, event.TicketID); err != nil {
		return err
	}

	return nil
}

func (h *Handler) TicketToRefund(ctx context.Context, event *entities.TicketBookingCanceled) error {
	if event.Price.Currency == "" {
		event.Price.Currency = "USD"
	}

	if err := h.spreadsheetsService.AppendRow(ctx, "tickets-to-refund", event.ToSpreadsheetTicketPayload()); err != nil {
		return err
	}
	return nil
}

func (h *Handler) BookPlaceInDeadNation(ctx context.Context, event *entities.BookingMade) error {
	log.FromContext(ctx).Info("Booking ticket in Dead Nation")

	show, err := h.showService.ShowByID(ctx, event.ShowId)
	if err != nil {
		return fmt.Errorf("failed to get show: %w", err)
	}

	err = h.deadNationAPI.BookInDeadNation(ctx, entities.DeadNationBooking{
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

func (h *Handler) IssueReceipt(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Issuing receipt")

	request := entities.IssueReceiptRequest{
		TicketID:       event.TicketID,
		Price:          event.Price,
		IdempotencyKey: event.Header.IdempotencyKey,
	}

	resp, err := h.receiptsService.IssueReceipt(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	if err = h.eventBus.Publish(ctx, entities.TicketReceiptIssued{
		Header:        event.Header,
		TicketID:      event.TicketID,
		ReceiptNumber: resp.ReceiptNumber,
		IssuedAt:      resp.IssuedAt,
	}); err != nil {
		logrus.Error("send message to TicketBookingConfirmed topic: ", err)
		return err
	}

	return nil
}

func (h *Handler) TicketEventHandlers() []cqrs.EventHandler {
	return []cqrs.EventHandler{
		cqrs.NewEventHandler("IssueReceipt", h.IssueReceipt),
		cqrs.NewEventHandler("TicketToPrint", h.TicketToPrint),
		cqrs.NewEventHandler("TicketToRefund", h.TicketToRefund),
		cqrs.NewEventHandler("SaveTicketInDB", h.SaveTicketInDB),
		cqrs.NewEventHandler("DeleteTicket", h.DeleteTicket),
		cqrs.NewEventHandler("StoreTicketContent", h.StoreTicketContent),
		cqrs.NewEventHandler("BookPlaceInDeadNation", h.BookPlaceInDeadNation),
	}
}

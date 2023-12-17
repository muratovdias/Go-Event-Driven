package v1

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
	"net/http"
	"tickets/internal/entities"
)

func (h *Handler) Tickets(c echo.Context) error {
	var tickets entities.TicketsStatusRequest
	if err := c.Bind(&tickets); err != nil {
		return err
	}
	// iterate over tickets
	for _, ticket := range tickets.Tickets {
		// add Header
		ticket.Header = entities.NewEventHeader()

		switch ticket.Status {
		case "confirmed":
			// publish message
			err := h.publisher.Publish(context.Background(), entities.TicketBookingConfirmed{
				Header:        ticket.Header,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			})
			if err != nil {
				h.watermillLogger.Error("send message to issue-receipts topic", err, watermill.LogFields{})
			}
		case "canceled":
			// publish message
			err := h.publisher.Publish(context.Background(), entities.TicketBookingCanceled{
				Header:        ticket.Header,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			})
			if err != nil {
				h.watermillLogger.Error("send message to issue-receipts topic", err, watermill.LogFields{})
			}
		}
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

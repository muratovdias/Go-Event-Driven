package v1

import (
	"context"
	"fmt"
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
		// get idempotency key
		idempotencyKey := c.Request().Header.Get("Idempotency-Key")
		if idempotencyKey == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Idempotency-Key header is required")
		}

		// add Header
		ticket.Header = entities.NewEventHeader(idempotencyKey + ticket.TicketID)

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
				h.watermillLogger.Error("send message to TicketBookingConfirmed topic", err, watermill.LogFields{})
			}

			err = h.publisher.Publish(context.Background(), entities.TicketPrinted{
				Header:   ticket.Header,
				TicketID: ticket.TicketID,
				FileName: fmt.Sprintf("%s-ticket.html", ticket.TicketID),
			})
			if err != nil {
				h.watermillLogger.Error("send message to TicketPrinted topic", err, watermill.LogFields{})
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

func (h *Handler) TicketsList(c echo.Context) error {
	tickets, err := h.service.TicketList(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tickets)
}

func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

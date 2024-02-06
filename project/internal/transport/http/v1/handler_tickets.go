package v1

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
			err := h.eventPublisher.Publish(c.Request().Context(), entities.TicketBookingConfirmed{
				Header:        ticket.Header,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			})
			if err != nil {
				logrus.Error("send message to TicketBookingConfirmed topic: ", err)
			}

			err = h.eventPublisher.Publish(c.Request().Context(), entities.TicketPrinted{
				Header:   ticket.Header,
				TicketID: ticket.TicketID,
				FileName: fmt.Sprintf("%s-ticket.html", ticket.TicketID),
			})
			if err != nil {
				logrus.Error("send message to TicketPrinted topic: ", err)
			}

		case "canceled":
			// publish message
			err := h.eventPublisher.Publish(c.Request().Context(), entities.TicketBookingCanceled{
				Header:        ticket.Header,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			})
			if err != nil {
				logrus.Error("send message to TicketBookingCanceled topic: ", err)
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

func (h *Handler) RefundTicket(c echo.Context) error {
	ticketID := c.Param("ticket_id")

	event := entities.RefundTicket{
		Header:   entities.NewEventHeader(uuid.NewString()),
		TicketID: ticketID,
	}

	if err := h.commandPublisher.Send(c.Request().Context(), event); err != nil {
		return fmt.Errorf("failed to send RefundTicket command: %w", err)
	}

	return c.NoContent(http.StatusAccepted)
}

func (h *Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

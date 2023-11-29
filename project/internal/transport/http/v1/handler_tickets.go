package v1

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
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
		// make event payload
		eventPayload, err := json.Marshal(ticket)
		if err != nil {
			return err
		}
		// create message
		msg := message.NewMessage(watermill.NewUUID(), eventPayload)
		msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))

		switch ticket.Status {
		case "confirmed":
			msg.Metadata.Set("type", "TicketBookingConfirmed")
			// publish message
			err = h.publisher.Publish("TicketBookingConfirmed", msg)
			if err != nil {
				h.watermillLogger.Error("send message to issue-receipts topic", err, watermill.LogFields{})
			}
		case "canceled":
			msg.Metadata.Set("type", "TicketBookingCanceled")
			// publish message
			err = h.publisher.Publish("TicketBookingCanceled", msg)
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

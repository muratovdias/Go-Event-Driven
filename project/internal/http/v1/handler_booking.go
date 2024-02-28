package v1

import (
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
	"net/http"
	"tickets/internal/entities"
	booking2 "tickets/internal/repository/booking"
)

func (h *Handler) BookTicket(c echo.Context) error {
	var booking entities.Booking

	if err := c.Bind(&booking); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	bookingID, err := h.service.BookTicket(c.Request().Context(), booking)
	if err != nil {
		if errors.As(err, &booking2.NotEnoughSeatsAvailableError{}) {
			h.watermillLogger.Error("", err, watermill.LogFields{"error": err.Error()})
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		h.watermillLogger.Error("", err, watermill.LogFields{"error": err.Error()})
		return echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"booking_id": bookingID,
	})
}

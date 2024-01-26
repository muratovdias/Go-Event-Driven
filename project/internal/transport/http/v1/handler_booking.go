package v1

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"tickets/internal/entities"
)

func (h *Handler) BookTicket(c echo.Context) error {
	var booking entities.Booking

	if err := c.Bind(&booking); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	bookingID, err := h.service.BookTicket(c.Request().Context(), booking)
	if err != nil {
		if errors.Is(err, fmt.Errorf("not enough seats available")) {
			return echo.NewHTTPError(http.StatusBadRequest, "not enough seats available")
		}
		return echo.NewHTTPError(http.StatusBadRequest, "not enough seats available")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"booking_id": bookingID,
	})
}

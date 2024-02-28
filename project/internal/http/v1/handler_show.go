package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tickets/internal/entities"
)

func (h *Handler) NewShow(c echo.Context) error {
	var show entities.Show

	if err := c.Bind(&show); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	showID, err := h.service.NewShow(c.Request().Context(), show)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"show_id": showID,
	})
}

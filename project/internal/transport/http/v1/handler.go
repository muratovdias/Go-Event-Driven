package v1

import (
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	publisher       publisherI
	watermillLogger loggerI
}

func NewHandler(publisher publisherI, watermillLogger loggerI) *Handler {
	return &Handler{publisher: publisher, watermillLogger: watermillLogger}
}

func (h *Handler) SetRoutes() *echo.Echo {
	router := commonHTTP.NewEcho()

	router.POST("/tickets-status", h.Tickets)
	router.GET("/health", h.Health)

	return router
}

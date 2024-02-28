package v1

import (
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SetRoutes() *echo.Echo {
	router := commonHTTP.NewEcho()

	router.POST("/tickets-status", h.Tickets)
	router.GET("/tickets", h.TicketsList)
	router.GET("/health", h.Health)
	router.POST("/shows", h.NewShow)
	router.POST("/book-tickets", h.BookTicket)
	router.PUT("/ticket-refund/:ticket_id", h.RefundTicket)

	return router
}

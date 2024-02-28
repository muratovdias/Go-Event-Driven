package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"tickets/internal/entities"
	"tickets/internal/repository/booking"
	"tickets/internal/repository/readModel"
	"tickets/internal/repository/show"
	"tickets/internal/repository/ticket"
)

type Ticket interface {
	SaveTicket(ctx context.Context, confirmed entities.TicketBookingConfirmed) error
	DeleteTicket(ctx context.Context, ticketID string) error
	TicketList(ctx context.Context) ([]entities.TicketList, error)
}

type Show interface {
	NewShow(ctx context.Context, show entities.Show) (string, error)
	ShowByID(ctx context.Context, showId uuid.UUID) (entities.Show, error)
}

type Booking interface {
	BookTicket(ctx context.Context, booking entities.Booking) (string, error)
}

type Ops interface {
	AllReservations() ([]entities.OpsBooking, error)
	ReservationReadModel(ctx context.Context, bookingID string) (entities.OpsBooking, error)
	OnBookingMade(ctx context.Context, bookingMade *entities.BookingMade) error
	OnTicketBookingConfirmed(ctx context.Context, event *entities.TicketBookingConfirmed) error
	OnTicketRefunded(ctx context.Context, event *entities.RefundTicket) error
	OnTicketPrinted(ctx context.Context, event *entities.TicketPrinted) error
}

type Repository struct {
	Ticket  Ticket
	Show    Show
	Booking Booking
	Ops     Ops
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Ticket:  ticket.NewRepo(db),
		Show:    show.NewRepo(db),
		Booking: booking.NewRepo(db),
		Ops:     readModel.NewOpsBookingReadModel(db),
	}
}

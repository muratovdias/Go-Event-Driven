package booking

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/jmoiron/sqlx"
	"tickets/internal/broker/event"
	"tickets/internal/broker/outbox"
	"tickets/internal/entities"
)

type Repo struct {
	db     *sqlx.DB
	logger watermill.LoggerAdapter
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db:     db,
		logger: watermill.NewStdLogger(true, true),
	}
}

type NotEnoughSeatsAvailableError struct {
	Available int
	Booked    int
}

func (e NotEnoughSeatsAvailableError) Error() string {
	return fmt.Sprintf("not enough seats available, available: %d, booked: %d", e.Available, e.Booked)
}

func (r *Repo) BookTicket(ctx context.Context, booking entities.Booking) (string, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return "", err
	}

	row := r.db.QueryRowContext(ctx, compareBeforeBooking, booking.ShowID)
	if err != nil {
		return "", err
	}
	var available, booked int

	err = row.Scan(&available, &booked)
	if err != nil {
		return "", err
	}

	if (available - booked) < booking.NumberOfTickets {
		tx.Rollback()
		return "", NotEnoughSeatsAvailableError{
			available,
			booked,
		}
	}

	rows, err := tx.QueryContext(ctx, inserBooking, booking.BookingID, booking.ShowID, booking.NumberOfTickets, booking.CustomerEmail)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return "", err
		}

		return "", err
	}
	defer rows.Close()

	var bookingID string
	for rows.Next() {
		if err := rows.Scan(&bookingID); err != nil {
			return "", err
		}
	}

	outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("could not create event bus: %w", err)
	}

	eventBus, err := event.NewEventBus(outboxPublisher)
	if err != nil {
		return "", fmt.Errorf("could not create event bus for book ticket: %w", err)
	}

	err = eventBus.Publish(ctx, entities.BookingMade{
		Header:          entities.NewEventHeader(""),
		BookingID:       booking.BookingID,
		NumberOfTickets: booking.NumberOfTickets,
		CustomerEmail:   booking.CustomerEmail,
		ShowId:          booking.ShowID,
	})
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return bookingID, err
}

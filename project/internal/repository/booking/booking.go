package booking

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/jmoiron/sqlx"
	"tickets/internal/broker"
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

func (r *Repo) BookTicket(ctx context.Context, booking entities.Booking) (string, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return "", err
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

	//payload, err := json.Marshal(booking)
	//if err != nil {
	//	return "", err
	//}

	//msg := watermillMessage.NewMessage(uuid.NewString(), payload)
	//
	//if err := PublishInTx(tx, msg, r.logger); err != nil {
	//	rbErr := tx.Rollback()
	//	if rbErr != nil {
	//		return "", err
	//	}
	//
	//	return "", err
	//}

	outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("could not create event bus: %w", err)
	}

	eventBus, err := broker.NewEventBus(outboxPublisher)
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

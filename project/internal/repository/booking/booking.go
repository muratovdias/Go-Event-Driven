package booking

import (
	"context"
	"github.com/jmoiron/sqlx"
	"tickets/internal/entities"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) BookTicket(ctx context.Context, booking entities.Booking) (string, error) {
	rows, err := r.db.NamedQueryContext(ctx, inserBooking, booking)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var bookingID string
	for rows.Next() {
		if err := rows.Scan(&bookingID); err != nil {
			return "", err
		}
	}

	return bookingID, err
}

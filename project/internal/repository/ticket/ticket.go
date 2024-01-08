package ticket

import (
	"context"
	"github.com/jmoiron/sqlx"
	"tickets/internal/entities"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) SaveTicket(ctx context.Context, confirmed entities.TicketBookingConfirmed) error {
	_, err := r.db.ExecContext(ctx, saveTicket,
		confirmed.TicketID,
		confirmed.Price.Amount,
		confirmed.Price.Currency,
		confirmed.CustomerEmail,
	)

	return err
}

func (r *Repo) DeleteTicket(ctx context.Context, ticketID string) error {
	_, err := r.db.ExecContext(ctx, deleteTicket, ticketID)
	return err
}

func (r *Repo) GetByID(ctx context.Context, ticketID string) (entities.Ticket, error) {
	row := r.db.QueryRowContext(ctx, getTicketByID, ticketID)

	var ticket entities.Ticket

	err := row.Scan(
		&ticket.TicketID,
		&ticket.Price.Amount,
		&ticket.Price.Currency,
		&ticket.CustomerEmail,
	)

	return ticket, err
}

func (r *Repo) TicketList(ctx context.Context) ([]entities.TicketList, error) {
	rows, err := r.db.QueryContext(ctx, ticketList)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tickets []entities.TicketList

	for rows.Next() {
		var ticket entities.TicketList

		if err := rows.Scan(
			&ticket.TicketID,
			&ticket.Price.Amount,
			&ticket.Price.Currency,
			&ticket.CustomerEmail,
		); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

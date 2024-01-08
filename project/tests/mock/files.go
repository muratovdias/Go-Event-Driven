package mock

import (
	"context"
	"sync"
	"tickets/internal/entities"
)

type FilesMock struct {
	mock    sync.Mutex
	Tickets map[string]struct{}
}

func (f *FilesMock) StoreTicketContent(ctx context.Context, ticket entities.TicketBookingConfirmed) error {

	if _, ok := f.Tickets[ticket.TicketID]; !ok {
		f.Tickets[ticket.TicketID] = struct{}{}
	}

	return nil
}

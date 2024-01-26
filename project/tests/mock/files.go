package mock

import (
	"context"
	"fmt"
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

func (f *FilesMock) DownloadTicketContent(ctx context.Context, ticketID string) (struct{}, error) {
	f.mock.Lock()
	defer f.mock.Unlock()

	if f.Tickets == nil {
		f.Tickets = make(map[string]struct{})
	}

	fileContent, ok := f.Tickets[ticketID]
	if !ok {
		return fileContent, fmt.Errorf("file %s not found", ticketID)
	}

	return fileContent, nil
}

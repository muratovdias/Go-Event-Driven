package mock

import (
	"context"
	"sync"
	"tickets/internal/entities"
)

type DeadNationClient struct {
	mx                 sync.Mutex
	DeadNationBookings []entities.DeadNationBooking
}

func (d *DeadNationClient) BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.DeadNationBookings = append(d.DeadNationBookings, request)

	return nil
}

package mock

import (
	"context"
	"sync"
	"tickets/internal/entities"
)

type PaymentsMock struct {
	lock    sync.Mutex
	Refunds []entities.PaymentRefund
}

func (c *PaymentsMock) PutRefundsWithResponse(ctx context.Context, refundPayment entities.PaymentRefund) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Refunds = append(c.Refunds, refundPayment)

	return nil
}

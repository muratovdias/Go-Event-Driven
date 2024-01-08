package db

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/stretchr/testify/assert"
	"testing"
	"tickets/internal/entities"
)

func testTicketBookingConfirmed(ticketID string) entities.TicketBookingConfirmed {
	return entities.TicketBookingConfirmed{
		TicketID: ticketID,
		Price: entities.Money{
			Amount:   "123",
			Currency: "USD",
		},
		CustomerEmail: "test",
	}
}

func TestTicket_Save_idempotency(t *testing.T) {
	ctx := context.Background()
	uuid := watermill.NewUUID()

	ticket := testTicketBookingConfirmed(uuid)

	for i := 0; i < 2; i++ {
		err := ticketRepo.SaveTicket(ctx, ticket)
		assert.NoError(t, err)

		got, err := ticketRepo.GetByID(ctx, ticket.TicketID)

		assert.NoError(t, err)
		assert.Equal(t, ticket.TicketID, got.TicketID)
		//assert.Equal(t, ticket.Price.Amount, got.Price.Amount)
		assert.Equal(t, ticket.Price.Currency, got.Price.Currency)
		assert.Equal(t, ticket.CustomerEmail, got.CustomerEmail)
	}
}

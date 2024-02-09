package payment

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/payments"
	"github.com/sirupsen/logrus"
	"net/http"
	"tickets/internal/entities"
)

type Client struct {
	clients *clients.Clients
}

func NewPaymentClient(clients *clients.Clients) *Client {
	if clients == nil {
		panic("NewFilesApiClient: clients is nil")
	}

	return &Client{clients: clients}
}

func (c *Client) PutRefundsWithResponse(ctx context.Context, command entities.RefundTicket) error {
	body := payments.PutRefundsJSONRequestBody{
		DeduplicationId:  &command.Header.IdempotencyKey,
		PaymentReference: command.TicketID,
		Reason:           "customer requested refund",
	}

	response, err := c.clients.Payments.PutRefundsWithResponse(ctx, body)
	if err != nil {
		logrus.Errorf("PutRefundsWithResponse: %v", err)
		return fmt.Errorf("failed to post refund for payment %s: %w", command.TicketID, err)
	}

	if response.StatusCode() != http.StatusOK {
		logrus.Infof("PutRefundsWithResponse status code: %d", response.StatusCode())
		return fmt.Errorf("unexpected for /payments-api/refunds status code: %d", response.StatusCode())
	}

	return nil
}

package receipts

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/receipts"
	"github.com/sirupsen/logrus"
	"net/http"
	"tickets/internal/entities"
)

type Client struct {
	clients *clients.Clients
}

func NewReceiptsClient(clients *clients.Clients) *Client {
	return &Client{
		clients: clients,
	}
}

func (c *Client) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	body := receipts.PutReceiptsJSONRequestBody{
		TicketId:       request.TicketID,
		IdempotencyKey: &request.IdempotencyKey,
		Price: receipts.Money{
			MoneyAmount:   request.Price.Amount,
			MoneyCurrency: request.Price.Currency,
		},
	}

	resp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, body)
	if err != nil {
		return entities.IssueReceiptResponse{}, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		// receipt already exists
		return entities.IssueReceiptResponse{
			ReceiptNumber: resp.JSON200.Number,
			IssuedAt:      resp.JSON200.IssuedAt,
		}, nil
	case http.StatusCreated:
		// receipt was created
		return entities.IssueReceiptResponse{
			ReceiptNumber: resp.JSON201.Number,
			IssuedAt:      resp.JSON201.IssuedAt,
		}, nil
	default:
		return entities.IssueReceiptResponse{}, fmt.Errorf("unexpected status code for POST receipts-api/receipts: %d", resp.StatusCode())
	}
}

func (c *Client) PutVoidReceiptWithResponse(ctx context.Context, command entities.VoidReceipt) error {
	body := receipts.PutVoidReceiptJSONRequestBody{
		Reason:       "customer requested refund",
		TicketId:     command.TicketID,
		IdempotentId: &command.IdempotencyKey,
	}

	response, err := c.clients.Receipts.PutVoidReceiptWithResponse(ctx, body)
	if err != nil {
		logrus.Errorf("PutVoidReceiptWithResponse: %v", err)
		return err
	}

	if response.StatusCode() != 200 {
		logrus.Infof("PutRefundsWithResponse status code: %d", response.StatusCode())
		return fmt.Errorf("unexpected status code for POST receipts-api/receipts/void: %d", response.StatusCode())
	}

	return nil
}

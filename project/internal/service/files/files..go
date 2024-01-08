package files

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"net/http"
	"tickets/internal/entities"
)

type Client struct {
	client *clients.Clients
}

func NewClient(client *clients.Clients) *Client {
	return &Client{client: client}
}

func (c *Client) StoreTicketContent(ctx context.Context, ticket entities.TicketBookingConfirmed) error {
	ticketContent := fmt.Sprintf(`<html>
			<head>
				<title>Ticket</title>
			</head>
			<body>
				<h1>Ticket %s</h1>
				<p>Price: %s %s</p>	
			</body>
		</html>`, ticket.TicketID, ticket.Price.Amount, ticket.Price.Currency)

	fileName := fmt.Sprintf("%s-ticket.html", ticket.TicketID)

	response, err := c.client.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, fileName, ticketContent)
	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusConflict {
		log.FromContext(ctx).Info("file %s already exists", ticket.TicketID)
		return nil
	}

	return nil
}

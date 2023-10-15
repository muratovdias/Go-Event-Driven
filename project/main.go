package main

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/receipts"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/spreadsheets"
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	log.Init(logrus.InfoLevel)
	watermillLogger := log.NewWatermill(logrus.NewEntry(logrus.StandardLogger()))

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, watermillLogger)
	if err != nil {
		watermillLogger.Error("creating new redis stream publisher", err, watermill.LogFields{})
		panic(err)
	}

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "receipts",
	}, watermillLogger)
	if err != nil {
		watermillLogger.Error("creating new redis stream subscriber", err, watermill.LogFields{})
		panic(err)
	}

	appendToTracker, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "spreadsheet",
	}, watermillLogger)
	if err != nil {
		watermillLogger.Error("creating new redis stream subscriber", err, watermill.LogFields{})
		panic(err)
	}

	client, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	receiptsClient := NewReceiptsClient(client, issueReceiptSub)
	spreadsheetsClient := NewSpreadsheetsClient(client, appendToTracker)

	e := commonHTTP.NewEcho()

	e.POST("/tickets-confirmation", func(c echo.Context) error {
		var request TicketsConfirmationRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		go spreadsheetsClient.ProcessingMessages()
		go receiptsClient.ProcessingMessages()

		for _, ticketID := range request.Tickets {
			msg := message.NewMessage(watermill.NewUUID(), []byte(ticketID))
			err := publisher.Publish("issue-receipt", msg)
			if err != nil {
				watermillLogger.Error("send message to issue-receipts topic", err, watermill.LogFields{})
			}
			err = publisher.Publish("append-to-tracker", msg)
			if err != nil {
				watermillLogger.Error("send message to append-to-tracker topic", err, watermill.LogFields{})
			}
		}

		return c.NoContent(http.StatusOK)
	})

	watermillLogger.Info("Server starting...", watermill.LogFields{})

	err = e.Start(":8080")
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

type TicketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

type ReceiptsClient struct {
	clients         *clients.Clients
	issueReceiptSub *redisstream.Subscriber
}

func NewReceiptsClient(clients *clients.Clients, subscriber *redisstream.Subscriber) ReceiptsClient {
	return ReceiptsClient{
		clients:         clients,
		issueReceiptSub: subscriber,
	}
}

func (c ReceiptsClient) IssueReceipt(ctx context.Context, ticketID string) error {
	body := receipts.PutReceiptsJSONRequestBody{
		TicketId: ticketID,
	}

	receiptsResp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if receiptsResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", receiptsResp.StatusCode())
	}

	return nil
}

func (c ReceiptsClient) ProcessingMessages() {
	messages, err := c.issueReceiptSub.Subscribe(context.Background(), "issue-receipt")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		ticketID := string(msg.Payload)
		if err := c.IssueReceipt(context.Background(), ticketID); err != nil {
			msg.Nack()
		} else {
			fmt.Printf("processed ticket: %s\n", ticketID)
			msg.Ack()
		}
	}
}

type SpreadsheetsClient struct {
	clients         *clients.Clients
	appendToTracker *redisstream.Subscriber
}

func NewSpreadsheetsClient(clients *clients.Clients, subscriber *redisstream.Subscriber) SpreadsheetsClient {
	return SpreadsheetsClient{
		clients:         clients,
		appendToTracker: subscriber,
	}
}

func (c SpreadsheetsClient) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	request := spreadsheets.PostSheetsSheetRowsJSONRequestBody{
		Columns: row,
	}

	sheetsResp, err := c.clients.Spreadsheets.PostSheetsSheetRowsWithResponse(ctx, spreadsheetName, request)
	if err != nil {
		return err
	}
	if sheetsResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", sheetsResp.StatusCode())
	}

	return nil
}

func (c SpreadsheetsClient) ProcessingMessages() {
	messages, err := c.appendToTracker.Subscribe(context.Background(), "append-to-tracker")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		ticketID := string(msg.Payload)
		if err := c.AppendRow(context.Background(), "tickets-to-print", []string{ticketID}); err != nil {
			msg.Nack()
		} else {
			fmt.Printf("processed ticket: %s\n", ticketID)
			msg.Ack()
		}
	}
}

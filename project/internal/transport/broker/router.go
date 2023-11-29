package broker

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/redis/go-redis/v9"
	"tickets/internal/entities"
	"time"
)

const (
	issueReceipt    = "issue-receipt"
	appendToTracker = "append-to-tracker"
	cancelTicket    = "cancel-ticket"

	malformedMessage = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"
)

type broker struct {
	service         serviceI
	rdb             *redis.Client
	watermillLogger watermill.LoggerAdapter
	router          *message.Router
	subscribers     map[string]*redisstream.Subscriber
}

func NewWatermillRouter(service serviceI, rdb *redis.Client, watermillLogger watermill.LoggerAdapter) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	broker := &broker{
		service:         service,
		rdb:             rdb,
		watermillLogger: watermillLogger,
		router:          router,
		subscribers:     make(map[string]*redisstream.Subscriber),
	}

	// initialize broker subscribers
	broker.initSubscribers()

	// set broker handlers
	broker.setHandler()

	// set middlewares
	broker.setMiddlewares()

	return router
}

func (b *broker) initSubscribers() {
	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        b.rdb,
		ConsumerGroup: "issue-receipt",
	}, b.watermillLogger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        b.rdb,
		ConsumerGroup: "append-to-tracker",
	}, b.watermillLogger)
	if err != nil {
		panic(err)
	}

	cancelTicketSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        b.rdb,
		ConsumerGroup: "cancel-ticket",
	}, b.watermillLogger)
	if err != nil {
		panic(err)
	}

	b.subscribers[issueReceipt] = issueReceiptSub
	b.subscribers[appendToTracker] = appendToTrackerSub
	b.subscribers[cancelTicket] = cancelTicketSub
}

func (b *broker) setHandler() {
	// Receipts router
	b.router.AddNoPublisherHandler(
		"receipts",
		"TicketBookingConfirmed",
		b.subscribers[issueReceipt],
		func(msg *message.Message) error {
			// Fixing an incorrect message type
			get := msg.Metadata.Get("type")
			if get != "TicketBookingConfirmed" {
				return nil
			}
			// Fixing a malformed JSON message
			if msg.UUID == malformedMessage {
				msg.Ack()
				return nil
			}
			var eventData entities.Ticket
			if err := json.Unmarshal(msg.Payload, &eventData); err != nil {
				return err
			}

			if eventData.Price.Currency == "" {
				eventData.Price.Currency = "USD"
			}

			if _, err := b.service.IssueReceipt(msg.Context(), eventData.ToIssueReceiptPayload()); err != nil {
				return err
			}
			return nil
		},
	)
	// Spreadsheet router
	b.router.AddNoPublisherHandler(
		"spreadsheet",
		"TicketBookingConfirmed",
		b.subscribers[appendToTracker],
		func(msg *message.Message) error {
			// Fixing an incorrect message type
			get := msg.Metadata.Get("type")
			if get != "TicketBookingConfirmed" {
				return nil
			}
			// Fixing a malformed JSON message
			if msg.UUID == malformedMessage {
				return nil
			}

			var eventData entities.Ticket
			if err := json.Unmarshal(msg.Payload, &eventData); err != nil {
				return err
			}

			if eventData.Price.Currency == "" {
				eventData.Price.Currency = "USD"
			}

			// add ticket
			if err := b.service.AppendRow(msg.Context(), "tickets-to-print", eventData.ToSpreadsheetTicketPayload()); err != nil {
				return err
			}
			return nil
		},
	)
	// Refund ticket router
	b.router.AddNoPublisherHandler(
		"refund ticket",
		"TicketBookingCanceled",
		b.subscribers[cancelTicket],
		func(msg *message.Message) error {
			// Fixing an incorrect message type
			get := msg.Metadata.Get("type")
			if get != "TicketBookingCanceled" {
				return nil
			}
			// Fixing a malformed JSON message
			if msg.UUID == malformedMessage {
				return nil
			}

			var eventData entities.Ticket
			if err := json.Unmarshal(msg.Payload, &eventData); err != nil {
				return err
			}

			if eventData.Price.Currency == "" {
				eventData.Price.Currency = "USD"
			}

			if err := b.service.AppendRow(msg.Context(), "tickets-to-refund", eventData.ToSpreadsheetTicketPayload()); err != nil {
				return err
			}
			return nil
		},
	)
}

func (b *broker) setMiddlewares() {
	// Retry failed messages
	retryMiddleware := middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          b.watermillLogger,
	}
	b.router.AddMiddleware(middleware.Recoverer, middleware.CorrelationID, PropagateCorrelationID, LoggingMiddleware, retryMiddleware.Middleware)
}

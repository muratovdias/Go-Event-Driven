package broker

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"tickets/internal/broker/outbox"
	"time"
)

type broker struct {
	eventHandlers   *eventHandlers
	watermillLogger watermill.LoggerAdapter
	router          *message.Router
	eventProcessor  *cqrs.EventProcessor
}

func NewWatermillRouter(service serviceI,
	postgresSubscriber message.Subscriber,
	publisher message.Publisher,
	eventProcessorConfig cqrs.EventProcessorConfig,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	// validate
	if service == nil {
		panic("missing service")
	}
	if postgresSubscriber == nil {
		panic("missing postgresSubscriber")
	}
	if publisher == nil {
		panic("missing publisher")
	}

	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	outbox.AddForwarderHandler(postgresSubscriber, publisher, router, watermillLogger)

	broker := &broker{
		watermillLogger: watermillLogger,
		router:          router,
	}

	// initialize event handlers
	broker.eventHandlers = newEventHandlers(service)

	// initialize broker subscribers
	broker.eventProcessor, err = cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	// set broker handlers
	broker.setEventHandler()

	// set middlewares
	broker.setMiddlewares()

	return router
}

func (b *broker) setEventHandler() {
	err := b.eventProcessor.AddHandlers(
		b.eventHandlers.ticketHandler.ticketHandlers()...,
	)
	if err != nil {
		panic(err)
	}
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
	b.router.AddMiddleware(middleware.Recoverer, PropagateCorrelationID, middleware.CorrelationID, LoggingMiddleware, retryMiddleware.Middleware)
}

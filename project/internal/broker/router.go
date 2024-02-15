package broker

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"tickets/internal/broker/outbox"
	"time"
)

type broker struct {
	eventHandlers    *eventHandlers
	commandHandlers  *commandHandlers
	watermillLogger  watermill.LoggerAdapter
	router           *message.Router
	eventProcessor   *cqrs.EventProcessor
	eventPublisher   eventPublisher
	commandProcessor *cqrs.CommandProcessor
}

func NewWatermillRouter(service serviceI,
	postgresSubscriber message.Subscriber,
	publisher message.Publisher,
	eventPublisher eventPublisher,
	eventProcessorConfig cqrs.EventProcessorConfig,
	commandProcessorConfig cqrs.CommandProcessorConfig,
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

	if eventPublisher == nil {
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
	broker.eventHandlers = newEventHandlers(service, eventPublisher)

	// initialize command handlers
	broker.commandHandlers = newCommandHandlers(service, eventPublisher)

	// initialize event subscriber
	broker.eventProcessor, err = cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(fmt.Errorf("initialize event subscriber failed: %w", err))
	}

	// initialize command subscriber
	broker.commandProcessor, err = cqrs.NewCommandProcessorWithConfig(router, commandProcessorConfig)
	if err != nil {
		panic(fmt.Errorf("initialize command subscriber failed: %w", err))
	}

	// set event handlers
	broker.setEventHandlers()

	// set command handlers
	broker.setCommandHandlers()

	// set middlewares
	broker.setMiddlewares()

	return router
}

func (b *broker) setEventHandlers() {
	err := b.eventProcessor.AddHandlers(
		b.eventHandlers.ticketHandler.ticketEventHandlers()...,
	)
	if err != nil {
		panic(err)
	}
}

func (b *broker) setCommandHandlers() {
	err := b.commandProcessor.AddHandlers(
		b.commandHandlers.ticketHandler.ticketCommandHandler()...,
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

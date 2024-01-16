package app

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	broker2 "tickets/internal/broker"
	"tickets/internal/broker/event"
	"tickets/internal/broker/outbox"
	"tickets/internal/repository"
	"tickets/internal/service"
	v1 "tickets/internal/transport/http/v1"
	"time"
)

type App struct {
	watermillRouter *message.Router
	echoRouter      *echo.Echo
	rdb             *redis.Client
}

func Initialize(
	receiptsClient receiptsClient,
	spreadsheetsClient spreadsheetsClient,
	filesClient filesClient,
	deadNationClient deadNationClient,
	redisClient *redis.Client,
	db *sqlx.DB,
) *App {
	// logger init
	log.Init(logrus.InfoLevel)
	watermillLogger := log.NewWatermill(logrus.NewEntry(logrus.StandardLogger()))

	// publisher init
	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: redisClient,
	}, watermillLogger)
	if err != nil {
		watermillLogger.Error("creating new redis stream publisher", err, watermill.LogFields{})
		panic(err)
	}

	// publisher decorator
	publisherDecorator := &log.CorrelationPublisherDecorator{Publisher: publisher}

	// event bus init
	eventBus, err := broker2.NewEventBus(publisherDecorator)

	// repository init
	repo := repository.NewRepository(db)

	// service init
	serv := service.NewService(receiptsClient, spreadsheetsClient, filesClient, deadNationClient, repo)

	// handler init
	handler := v1.NewHandler(eventBus, serv, watermillLogger)

	// postgres subscriber
	postgresSubscriber := outbox.NewPostgresSubscriber(db, watermillLogger)

	// event processor config
	eventProcessorConfig := event.NewProcessorConfig(redisClient, watermillLogger)

	// broker router init
	brokerRouter := broker2.NewWatermillRouter(serv, postgresSubscriber, publisherDecorator,
		eventProcessorConfig, watermillLogger)

	// set http routes
	httpRouter := handler.SetRoutes()

	return &App{
		echoRouter:      httpRouter,
		watermillRouter: brokerRouter,
		rdb:             redisClient,
	}
}

func (a *App) Start() {

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.watermillRouter.Run(ctx)
	})

	g.Go(func() error {
		<-a.watermillRouter.Running()

		err := a.echoRouter.Start(":8080")
		if err != nil || errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		err := a.rdb.Close()
		if err != nil {
			return err
		}

		ctx, cancel = context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		err = a.echoRouter.Shutdown(ctx)
		if err != nil {
			return err
		}

		err = a.watermillRouter.Close()
		if err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

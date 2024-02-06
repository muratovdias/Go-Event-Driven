package app

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
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
	"tickets/internal/broker/command"
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
	paymentClient paymentClient,
	redisClient *redis.Client,
	db *sqlx.DB,
) *App {
	// logger init
	log.Init(logrus.InfoLevel)
	watermillLogger := log.NewWatermill(logrus.NewEntry(logrus.StandardLogger()))

	// publisher init
	publisher := broker2.NewRedisPublisher(redisClient, watermillLogger)

	// publisher decorator
	publisher = &log.CorrelationPublisherDecorator{Publisher: publisher}

	// event bus init
	eventBus, err := event.NewEventBus(publisher)
	if err != nil {
		panic(err)
	}

	// command bus init
	commandBus, err := command.NewCommandBus(publisher)
	if err != nil {
		panic(err)
	}

	// repository init
	repo := repository.NewRepository(db)

	// service init
	serv := service.NewService(receiptsClient, spreadsheetsClient, filesClient, deadNationClient, paymentClient, repo)

	// handler init
	handler := v1.NewHandler(eventBus, commandBus, serv, watermillLogger)

	// postgres subscriber
	postgresSubscriber := outbox.NewPostgresSubscriber(db, watermillLogger)

	// event processor config
	eventProcessorConfig := event.NewEventProcessorConfig(redisClient, watermillLogger)

	// command processor config
	commandProcessorConfig := command.NewCommandProcessorConfig(redisClient, watermillLogger)

	// broker router init
	brokerRouter := broker2.NewWatermillRouter(serv, postgresSubscriber, publisher,
		eventProcessorConfig, commandProcessorConfig, watermillLogger)

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

	//ctxWithCancel, cancel := context.WithCancel(ctx)

	g.Go(func() error {
		err := a.watermillRouter.Run(ctx)
		return err
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

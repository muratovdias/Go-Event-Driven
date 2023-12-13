package app

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"tickets/internal/service"
	"tickets/internal/transport/broker"
	v1 "tickets/internal/transport/http/v1"
	"time"
)

type App struct {
	watermillRouter *message.Router
	echoRouter      *echo.Echo
	rdb             *redis.Client
}

func Initialize(
	receiptClient receiptsClient,
	spreadsheetClient spreadsheetsClient,
	redisClient *redis.Client,
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
	eventBus, err := broker.NewEventBus(publisher)

	// handler init
	handler := v1.NewHandler(eventBus, watermillLogger)

	// service init
	serv := service.NewService(receiptClient, spreadsheetClient)

	// broker router init
	brokerRouter := broker.NewWatermillRouter(serv, redisClient, watermillLogger)

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

		ctx, _ = context.WithTimeout(ctx, time.Second*30)
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

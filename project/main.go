package main

import (
	"context"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/redis/go-redis/v9"
	_ "github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"tickets/internal/app"
	"tickets/internal/repository"
	"tickets/internal/service/deadnation"
	"tickets/internal/service/files"
	"tickets/internal/service/payment"
	"tickets/internal/service/receipts"
	"tickets/internal/service/spreadsheet"
)

func main() {
	// client init
	client, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	// API clients
	receiptsClient := receipts.NewReceiptsClient(client)
	spreadsheetsClient := spreadsheet.NewSpreadsheetsClient(client)
	filesClient := files.NewClient(client)
	deadNationClient := deadnation.NewDeadNationClient(client)
	paymentClient := payment.NewPaymentClient(client)

	db, err := repository.InitDB()
	if err != nil {
		panic(err)
	}

	// redis client init
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	app1 := app.Initialize(receiptsClient, spreadsheetsClient, filesClient, deadNationClient, paymentClient, rdb, db)
	app1.Start()
}

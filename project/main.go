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

	receiptClient := receipts.NewReceiptsClient(client)
	spreadsheetClient := spreadsheet.NewSpreadsheetsClient(client)

	// redis client init
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	app1 := app.Initialize(receiptClient, spreadsheetClient, rdb)
	app1.Start()
}

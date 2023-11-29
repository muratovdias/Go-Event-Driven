package tests_test

import (
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"testing"
	"tickets/internal/app"
	"tickets/tests/mock"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponent(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	defer rdb.Close()

	receiptClient := &mock.ReceiptMock{}
	spreadsheetClient := &mock.SpreadsheetsMock{}

	app1 := app.Initialize(receiptClient, spreadsheetClient, rdb)
	go app1.Start()

	waitForHttpServer(t)
}

func waitForHttpServer(t *testing.T) {
	t.Helper()

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8080/health")
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()

			if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
				return
			}
		},
		time.Second*10,
		time.Millisecond*50,
	)
}

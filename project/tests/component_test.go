package tests_test

import (
	"bytes"
	"encoding/json"
	"github.com/lithammer/shortuuid/v3"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"testing"
	"tickets/internal/app"
	"tickets/internal/entities"
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

	sendTicketsStatus(t, testTicketStatusRequest())

	assertReceiptForTicketIssued(t, receiptClient, testTicket("2"))
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

func sendTicketsStatus(t *testing.T, req entities.TicketsStatusRequest) {
	t.Helper()

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	correlationID := shortuuid.New()

	ticketIDs := make([]string, 0, len(req.Tickets))
	for _, ticket := range req.Tickets {
		ticketIDs = append(ticketIDs, ticket.TicketID)
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/tickets-status",
		bytes.NewBuffer(payload),
	)
	require.NoError(t, err)

	httpReq.Header.Set("Correlation-ID", correlationID)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	defer resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func assertReceiptForTicketIssued(t *testing.T, receiptsService *mock.ReceiptMock, ticket entities.Ticket) {
	assert.EventuallyWithT(
		t,
		func(collectT *assert.CollectT) {
			issuedReceipts := len(receiptsService.IssuedReceipts)
			t.Log("issued receipts", issuedReceipts)

			assert.Greater(collectT, issuedReceipts, 0, "no receipts issued")
		},
		10*time.Second,
		100*time.Millisecond,
	)

	var receipt entities.IssueReceiptRequest
	var ok bool
	for _, issuedReceipt := range receiptsService.IssuedReceipts {
		if issuedReceipt.TicketID != ticket.TicketID {
			continue
		}
		receipt = issuedReceipt
		ok = true
		break
	}
	require.Truef(t, ok, "receipt for ticket %s not found", ticket.TicketID)

	assert.Equal(t, ticket.TicketID, receipt.TicketID)
	assert.Equal(t, ticket.Price.Amount, receipt.Price.Amount)
	assert.Equal(t, ticket.Price.Currency, receipt.Price.Currency)
}

func testTicketStatusRequest() entities.TicketsStatusRequest {
	return entities.TicketsStatusRequest{
		Tickets: []entities.Ticket{
			testTicket("1"),
			testTicket("2"),
		},
	}
}

func testTicket(id string) entities.Ticket {
	return entities.Ticket{
		TicketID: id,
		Status:   "confirmed",
		Price: entities.Money{
			Amount:   "test",
			Currency: "test",
		},
	}
}

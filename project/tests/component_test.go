package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v3"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"testing"
	"tickets/internal/app"
	"tickets/internal/entities"
	"tickets/internal/repository"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rdb.Ping(ctx)

	db, err := repository.InitDB()
	if err != nil {
		panic(err)
	}

	receiptClient := &mock.ReceiptMock{IssuedReceipts: map[string]entities.IssueReceiptRequest{}}
	spreadsheetClient := &mock.SpreadsheetsMock{Rows: make(map[string][][]string)}
	filesClient := &mock.FilesMock{Tickets: make(map[string]struct{})}
	deadNationClient := &mock.DeadNationClient{DeadNationBookings: make([]entities.DeadNationBooking, 0)}
	paymentsService := &mock.PaymentsMock{}

	app1 := app.Initialize(receiptClient, spreadsheetClient, filesClient, deadNationClient, paymentsService, rdb, db)
	go app1.Start()

	waitForHttpServer(t)

	sendTicketsStatus(t, testTicketStatusRequest())

	assertReceiptForTicketIssued(t, receiptClient, testTicket("2", "confirmed"))
	assertTicketPrinted(t, filesClient, testTicket("2", "confirmed"))
	assertRowToSheetAdded(t, spreadsheetClient, testTicket("2", "confirmed"), "tickets-to-print")

}

func waitForHttpServer(t *testing.T) {
	t.Helper()
	<-time.After(time.Second * 2)

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8000/health")
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

	ticketsID := make([]string, 0, len(req.Tickets))
	for _, ticket := range req.Tickets {
		ticketsID = append(ticketsID, ticket.TicketID)
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8000/tickets-status",
		bytes.NewBuffer(payload),
	)
	require.NoError(t, err)

	httpReq.Header.Set("Correlation-ID", correlationID)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotency-Key", uuid.NewString())

	resp, err := http.DefaultClient.Do(httpReq)
	defer resp.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

/// asserts function START ///

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

func assertTicketPrinted(t *testing.T, filesApi *mock.FilesMock, ticket entities.Ticket) {
	assert.EventuallyWithT(
		t,
		func(collect *assert.CollectT) {
			content, err := filesApi.DownloadTicketContent(context.Background(), ticket.TicketID)

			require.NoError(t, err)
			assert.Equal(t, struct{}{}, content)
		},
		time.Second*10,
		time.Millisecond*100,
	)
}

func assertRowToSheetAdded(t *testing.T, spreadsheetsService *mock.SpreadsheetsMock, ticket entities.Ticket, sheetName string) bool {
	return assert.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			rows, ok := spreadsheetsService.Rows[sheetName]
			if !assert.True(t, ok, "sheet %s not found", sheetName) {
				return
			}

			allValues := []string{}

			for _, row := range rows {
				for _, col := range row {
					allValues = append(allValues, col)
				}
			}

			assert.Contains(t, allValues, ticket.TicketID, "ticket id not found in sheet %s", sheetName)
		},
		10*time.Second,
		100*time.Millisecond,
	)
}

/// asserts function END ///

func testTicketStatusRequest() entities.TicketsStatusRequest {
	return entities.TicketsStatusRequest{
		Tickets: []entities.Ticket{
			testTicket("1", "confirmed"),
			testTicket("2", "confirmed"),
		},
	}
}

func testTicket(id, status string) entities.Ticket {
	return entities.Ticket{
		TicketID: id,
		Status:   status,
		Price: entities.Money{
			Amount:   "test",
			Currency: "test",
		},
	}
}

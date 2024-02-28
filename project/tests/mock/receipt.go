package mock

import (
	"context"
	"sync"
	"tickets/internal/entities"
	"time"
)

type ReceiptMock struct {
	mock sync.Mutex

	IssuedReceipts map[string]entities.IssueReceiptRequest
	VoidedReceipts []entities.VoidReceipt
}

func (r *ReceiptMock) PutVoidReceiptWithResponse(ctx context.Context, request entities.VoidReceipt) error {
	r.mock.Lock()
	defer r.mock.Unlock()

	r.VoidedReceipts = append(r.VoidedReceipts, request)

	return nil
}

func (r *ReceiptMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	r.mock.Lock()
	defer r.mock.Unlock()

	r.IssuedReceipts[request.TicketID] = request

	return entities.IssueReceiptResponse{
		ReceiptNumber: "mocked-receipt-number",
		IssuedAt:      time.Now(),
	}, nil
}

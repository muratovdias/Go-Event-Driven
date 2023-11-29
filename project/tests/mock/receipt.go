package mock

import (
	"context"
	"sync"
	"tickets/internal/entities"
	"time"
)

type ReceiptMock struct {
	mock sync.Mutex

	IssuedReceipts []entities.IssueReceiptRequest
}

func (r *ReceiptMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	r.mock.Lock()
	defer r.mock.Unlock()

	r.IssuedReceipts = append(r.IssuedReceipts, request)

	return entities.IssueReceiptResponse{
		ReceiptNumber: "mocked-number",
		IssuedAt:      time.Now(),
	}, nil
}

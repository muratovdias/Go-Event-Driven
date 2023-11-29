package broker

import (
	"context"
	"tickets/internal/entities"
)

type serviceI interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}

package service

import (
	"context"
	"tickets/internal/entities"
)

type ReceiptsClient interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type SpreadsheetsClient interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}

type Service struct {
	ReceiptsClient
	SpreadsheetsClient
}

func NewService(receiptClient ReceiptsClient, spreadsheetClient SpreadsheetsClient) *Service {
	return &Service{
		ReceiptsClient:     receiptClient,
		SpreadsheetsClient: spreadsheetClient,
	}
}

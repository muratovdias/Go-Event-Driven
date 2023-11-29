package mock

import (
	"context"
	"sync"
	"tickets/internal/entities"
)

type SpreadsheetsMock struct {
	mock sync.Mutex

	IssuedReceipts []entities.IssueReceiptRequest
}

func (s *SpreadsheetsMock) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	s.mock.Lock()
	defer s.mock.Unlock()

	return nil
}

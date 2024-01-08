package mock

import (
	"context"
	"sync"
)

type SpreadsheetsMock struct {
	mock sync.Mutex

	Rows map[string][][]string
}

func (s *SpreadsheetsMock) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	s.mock.Lock()
	defer s.mock.Unlock()

	s.Rows[spreadsheetName] = append(s.Rows[spreadsheetName], row)

	return nil
}

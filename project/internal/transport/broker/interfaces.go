package broker

import (
	"tickets/internal/service"
)

type serviceI interface {
	service.ReceiptsClient
	service.SpreadsheetsClient
	service.FilesClient
	service.Ticket
}

package event

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"tickets/internal/entities"
)

type Handler struct {
	deadNationAPI       DeadNationAPI
	spreadsheetsService SpreadsheetsAPI
	receiptsService     ReceiptsService
	filesAPI            FilesAPI
	ticketService       TicketService
	showService         Show
	bookingService      Booking
	eventBus            *cqrs.EventBus
}

func NewHandler(
	deadNationAPI DeadNationAPI,
	spreadsheetsService SpreadsheetsAPI,
	receiptsService ReceiptsService,
	filesAPI FilesAPI,
	ticketService TicketService,
	showService Show,
	bookingService Booking,
	eventBus *cqrs.EventBus,
) Handler {
	if eventBus == nil {
		panic("missing eventBus")
	}
	if deadNationAPI == nil {
		panic("missing deadNationAPI")
	}
	if spreadsheetsService == nil {
		panic("missing spreadsheetsService")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	if filesAPI == nil {
		panic("missing filesAPI")
	}
	if eventBus == nil {
		panic("missing eventBus")
	}
	if ticketService == nil {
		panic("missing ticketService")
	}
	if showService == nil {
		panic("missing showService")
	}
	if bookingService == nil {
		panic("missing bookingService")
	}

	return Handler{
		deadNationAPI:       deadNationAPI,
		spreadsheetsService: spreadsheetsService,
		receiptsService:     receiptsService,
		filesAPI:            filesAPI,
		eventBus:            eventBus,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type FilesAPI interface {
	StoreTicketContent(ctx context.Context, ticket entities.TicketBookingConfirmed) error
}

type TicketService interface {
	SaveTicket(ctx context.Context, ticket entities.TicketBookingConfirmed) error
	DeleteTicket(ctx context.Context, ticketID string) error
	TicketList(ctx context.Context) ([]entities.TicketList, error)
}

type DeadNationAPI interface {
	BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error
}

type Show interface {
	NewShow(ctx context.Context, show entities.Show) (string, error)
	ShowByID(ctx context.Context, showId uuid.UUID) (entities.Show, error)
}

type Booking interface {
	BookTicket(ctx context.Context, booking entities.Booking) (string, error)
}

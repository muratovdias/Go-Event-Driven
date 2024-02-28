package command

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"tickets/internal/entities"
)

type Handler struct {
	eventBus *cqrs.EventBus

	receiptsServiceClient ReceiptsService
	paymentsServiceClient PaymentsService
}

func NewHandler(eventBus *cqrs.EventBus, receiptsServiceClient ReceiptsService, paymentsServiceClient PaymentsService) Handler {
	if eventBus == nil {
		panic("eventBus is required")
	}
	if receiptsServiceClient == nil {
		panic("receiptsServiceClient is required")
	}
	if paymentsServiceClient == nil {
		panic("paymentsServiceClient is required")
	}

	handler := Handler{
		eventBus:              eventBus,
		receiptsServiceClient: receiptsServiceClient,
		paymentsServiceClient: paymentsServiceClient,
	}

	return handler
}

func (h *Handler) TicketCommandHandler() []cqrs.CommandHandler {
	return []cqrs.CommandHandler{
		cqrs.NewCommandHandler("RefundTicket", h.RefundTicket),
	}
}

type ReceiptsService interface {
	PutVoidReceiptWithResponse(ctx context.Context, request entities.VoidReceipt) error
}

type PaymentsService interface {
	PutRefundsWithResponse(ctx context.Context, request entities.PaymentRefund) error
}

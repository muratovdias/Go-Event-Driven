package v1

type Handler struct {
	eventPublisher   eventPublisher
	commandPublisher commandSender
	service          serviceI
	watermillLogger  loggerI
}

func NewHandler(eventPublisher eventPublisher, commandPublisher commandSender, service serviceI, watermillLogger loggerI) *Handler {
	return &Handler{
		eventPublisher:   eventPublisher,
		commandPublisher: commandPublisher,
		service:          service,
		watermillLogger:  watermillLogger,
	}
}

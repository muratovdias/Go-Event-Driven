package v1

type Handler struct {
	eventPublisher   publisher
	commandPublisher sender
	service          serviceI
	watermillLogger  loggerI
}

func NewHandler(eventPublisher publisher, commandPublisher sender, service serviceI, watermillLogger loggerI) *Handler {
	return &Handler{
		eventPublisher:   eventPublisher,
		commandPublisher: commandPublisher,
		service:          service,
		watermillLogger:  watermillLogger,
	}
}

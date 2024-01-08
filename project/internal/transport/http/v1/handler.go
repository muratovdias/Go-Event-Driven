package v1

type Handler struct {
	publisher       publisherI
	service         serviceI
	watermillLogger loggerI
}

func NewHandler(publisher publisherI, service serviceI, watermillLogger loggerI) *Handler {
	return &Handler{
		publisher:       publisher,
		service:         service,
		watermillLogger: watermillLogger,
	}
}

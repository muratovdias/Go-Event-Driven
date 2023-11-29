package entities

type Ticket struct {
	Header        EventHeader `json:"header,omitempty"`
	TicketID      string      `json:"ticket_id,omitempty"`
	Status        string      `json:"status,omitempty"`
	CustomerEmail string      `json:"customer_email,omitempty"`
	Price         Money       `json:"price,omitempty"`
}

type TicketsStatusRequest struct {
	Tickets []Ticket `json:"tickets"`
}

func (t *Ticket) ToSpreadsheetTicketPayload() []string {
	return []string{t.TicketID, t.CustomerEmail, t.Price.Amount, t.Price.Currency}
}

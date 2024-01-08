package entities

type Ticket struct {
	Header        EventHeader `json:"header,omitempty"`
	TicketID      string      `json:"ticket_id,omitempty"`
	Status        string      `json:"status,omitempty"`
	CustomerEmail string      `json:"customer_email,omitempty"`
	Price         Money       `json:"price,omitempty"`
}

type TicketList struct {
	TicketID      string `json:"ticket_id,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`
	Price         Money  `json:"price,omitempty"`
}

type TicketsStatusRequest struct {
	Tickets []Ticket `json:"tickets"`
}

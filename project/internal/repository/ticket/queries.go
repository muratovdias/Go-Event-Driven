package ticket

const (
	saveTicket = `
INSERT INTO tickets (
  ticket_id, price_amount, price_currency, customer_email
) VALUES (
  $1, $2, $3, $4
)
ON CONFLICT DO NOTHING
`
	deleteTicket = `
DELETE
FROM tickets
WHERE ticket_id = $1
`

	getTicketByID = `
SELECT ticket_id, price_amount, price_currency, customer_email
FROM tickets
WHERE ticket_id = $1 LIMIT 1
`

	ticketList = `
SELECT ticket_id, price_amount, price_currency, customer_email
FROM tickets
`
)

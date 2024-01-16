package booking

const (
	inserBooking = `
INSERT INTO bookings (booking_id, show_id, number_of_tickets, customer_email)
VALUES ($1, $2, $3, $4)
RETURNING booking_id
`
)

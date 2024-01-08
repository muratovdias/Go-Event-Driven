package booking

const (
	inserBooking = `
INSERT INTO bookings (booking_id, show_id, number_of_tickets, customer_email)
VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email)
RETURNING booking_id
`
)

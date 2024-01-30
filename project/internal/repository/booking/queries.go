package booking

const (
	inserBooking = `
INSERT INTO bookings (booking_id, show_id, number_of_tickets, customer_email)
VALUES ($1, $2, $3, $4)
RETURNING booking_id
`

	compareBeforeBooking = `
SELECT
  s.number_of_tickets AS available_tickets,
  COALESCE(SUM(b.number_of_tickets), 0) AS booked_tickets
FROM
  shows s
LEFT JOIN
  bookings b ON s.show_id = b.show_id
WHERE
  s.show_id = $1
GROUP BY
  s.show_id, s.number_of_tickets;

`
)

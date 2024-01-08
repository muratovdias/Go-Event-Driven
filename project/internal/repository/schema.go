package repository

const (
	TicketsTable = `CREATE TABLE IF NOT EXISTS tickets (
	ticket_id UUID PRIMARY KEY,
	price_amount NUMERIC(10, 2) NOT NULL,
	price_currency CHAR(3) NOT NULL,
	customer_email VARCHAR(255) NOT NULL
);`

	ShowsTable = `CREATE TABLE IF NOT EXISTS shows(
    show_id UUID PRIMARY KEY,
    dead_nation_id UUID NOT NULL,
    number_of_tickets INTEGER NOT NULL,
    start_time TIMESTAMP NOT NULL,
    title VARCHAR NOT NULL,
    venue VARCHAR NOT NULL
);`

	BookingTable = `CREATE TABLE IF NOT EXISTS bookings(
	booking_id UUID PRIMARY KEY,
	show_id UUID NOT NULL,
	number_of_tickets INTEGER NOT NULL,
	customer_email VARCHAR NOT NULL 
);`
)

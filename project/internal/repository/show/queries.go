package show

const (
	insertShow = `
INSERT INTO shows (
  show_id, dead_nation_id, number_of_tickets, start_time, title, venue
) VALUES (
  :show_id, :dead_nation_id, :number_of_tickets, :start_time, :title, :venue
)
RETURNING show_id
`
)

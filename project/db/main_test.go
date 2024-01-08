package db

import (
	"os"
	"testing"
	"tickets/internal/repository/ticket"
)

var ticketRepo *ticket.Repo

func TestMain(m *testing.M) {
	db := getDB()

	ticketRepo = ticket.NewRepo(db)

	os.Exit(m.Run())
}

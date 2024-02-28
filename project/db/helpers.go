package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sync"
	"tickets/internal/repository"
)

var db *sqlx.DB
var getDbOnce sync.Once

func getDB() *sqlx.DB {
	getDbOnce.Do(func() {
		var err error

		db, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			log.Fatal(err)
		}
	})

	_, err := db.Exec(repository.SchemaPostgres)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

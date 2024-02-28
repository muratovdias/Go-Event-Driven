package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func InitDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sqlx.DB) error {
	if _, err := db.Exec(SchemaPostgres); err != nil {
		return err
	}

	return nil
}

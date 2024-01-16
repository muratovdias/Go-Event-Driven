package show

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"tickets/internal/entities"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) NewShow(ctx context.Context, show entities.Show) (string, error) {
	row, err := r.db.NamedQueryContext(ctx, insertShow, show)
	if err != nil {
		return "", err
	}
	defer row.Close()

	var showID string
	for row.Next() {
		err = row.Scan(&showID)
		if err != nil {
			return "", err
		}
	}

	return showID, err
}

func (r *Repo) ShowByID(ctx context.Context, showId uuid.UUID) (entities.Show, error) {
	var show entities.Show
	var err = r.db.GetContext(ctx, &show, `
		SELECT * 
		FROM shows
		WHERE show_id = $1
	`, showId)
	if err != nil {
		return entities.Show{}, fmt.Errorf("could not get show: %w", err)
	}

	return show, nil
}

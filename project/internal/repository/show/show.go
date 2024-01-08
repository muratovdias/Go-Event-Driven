package show

import (
	"context"
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

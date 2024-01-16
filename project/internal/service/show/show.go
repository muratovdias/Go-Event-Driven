package show

import (
	"context"
	"github.com/google/uuid"
	"tickets/internal/entities"
	"tickets/internal/repository"
)

type Service struct {
	repo repository.Show
}

func NewService(repo repository.Show) *Service {
	return &Service{repo: repo}
}

func (s *Service) NewShow(ctx context.Context, show entities.Show) (string, error) {
	show.ShowID = uuid.NewString()
	return s.repo.NewShow(ctx, show)
}

func (s *Service) ShowByID(ctx context.Context, showId uuid.UUID) (entities.Show, error) {
	return s.repo.ShowByID(ctx, showId)
}

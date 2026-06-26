package audit

import (
	"context"

	"github.com/example/ag-directory-manager/internal/models"
	"github.com/example/ag-directory-manager/internal/repositories"
)

type Service struct {
	repo *repositories.AuditRepository
}

func NewService(repo *repositories.AuditRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, log models.AuditLog) {
	_ = s.repo.Insert(ctx, log)
}

func (s *Service) List(ctx context.Context, limit int) ([]models.AuditLog, error) {
	return s.repo.List(ctx, limit)
}

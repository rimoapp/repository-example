package service

import (
	"context"
	"time"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type OrganizationService struct {
	Repo repository.OrganizationRepository
	BaseGenericService[*model.Organization]
}

func NewOrganizationService(repo repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{Repo: repo}
}

func (s *OrganizationService) Create(ctx context.Context, object *model.Organization) (string, error) {
	now := time.Now()
	object.CreatedAt = now
	object.UpdatedAt = now
	return s.Repo.Create(ctx, object)
}

func (s *OrganizationService) List(ctx context.Context, opts *model.OrganizationListOption) ([]*model.Organization, error) {
	return s.Repo.List(ctx, opts)
}

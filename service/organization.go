package service

import (
	"context"
	"time"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type OrganizationService interface {
	AbstractGenericService[*model.Organization, *model.OrganizationListOption]
}

type organizationService struct {
	repo repository.OrganizationRepository
	baseGenericService[*model.Organization, *model.OrganizationListOption]
}

func NewOrganizationService(repo repository.OrganizationRepository) OrganizationService {
	base := newGenericService(repo)
	return &organizationService{repo: repo, baseGenericService: base}
}

func (s *organizationService) Create(ctx context.Context, object *model.Organization) (string, error) {
	now := time.Now()
	object.CreatedAt = now
	object.UpdatedAt = now
	return s.repo.Create(ctx, object)
}

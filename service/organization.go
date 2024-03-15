package service

import (
	"context"
	"time"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type OrganizationService struct {
	Repo repository.OrganizationRepository
	BaseGenericService[*model.Organization, *model.OrganizationListOption]
	TeamService *TeamService
}

func NewOrganizationService(repo repository.OrganizationRepository, teamService *TeamService) *OrganizationService {
	base := NewGenericService(repo)
	return &OrganizationService{Repo: repo, BaseGenericService: base, TeamService: teamService}
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

func (s *OrganizationService) GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error) {
	org, err := s.Repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if opts != nil && opts.IncludeTeams {
		teams, err := s.TeamService.List(ctx, &model.TeamListOption{OrganizationID: id})
		if err != nil {
			return nil, err
		}
		org.Teams = teams
	}
	return org, nil
}

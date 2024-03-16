package service

import (
	"context"
	"time"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type OrganizationService interface {
	AbstractGenericService[*model.Organization, *model.OrganizationListOption]
	GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error)
}

type organizationService struct {
	repo repository.OrganizationRepository
	baseGenericService[*model.Organization, *model.OrganizationListOption]
	teamService TeamService
}

func NewOrganizationService(repo repository.OrganizationRepository, teamService TeamService) OrganizationService {
	base := newGenericService(repo)
	return &organizationService{repo: repo, baseGenericService: base, teamService: teamService}
}

func (s *organizationService) Create(ctx context.Context, object *model.Organization) (string, error) {
	now := time.Now()
	object.CreatedAt = now
	object.UpdatedAt = now
	return s.repo.Create(ctx, object)
}

func (s *organizationService) GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error) {
	org, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if opts != nil && opts.IncludeTeams {
		teams, err := s.teamService.List(ctx, &model.TeamListOption{OrganizationID: id})
		if err != nil {
			return nil, err
		}
		org.Teams = teams
	}
	return org, nil
}

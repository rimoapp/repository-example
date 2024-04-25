package usecase

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
)

type OrganizationUseCase interface {
	AbstractGenericUseCase[*model.Organization, *model.OrganizationListOption]
	GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error)
}

type organizationUseCase struct {
	baseGenericUseCase[*model.Organization, *model.OrganizationListOption]
	svc         service.OrganizationService
	teamService service.TeamService
}

func NewOrganizationUseCase(svc service.OrganizationService, teamService service.TeamService) OrganizationUseCase {
	base := newGenericUseCase(svc)
	return &organizationUseCase{baseGenericUseCase: base, svc: svc, teamService: teamService}
}

var _ OrganizationUseCase = &organizationUseCase{}

func (uc *organizationUseCase) GetWithOption(ctx context.Context, id string, opts *model.GetOrganizationOption) (*model.Organization, error) {
	org, err := uc.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if opts != nil && opts.IncludeTeams {
		teams, err := uc.teamService.List(ctx, &model.TeamListOption{OrganizationID: id})
		if err != nil {
			return nil, err
		}
		org.Teams = teams
	}
	return org, nil
}

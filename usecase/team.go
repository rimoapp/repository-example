package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
)

type TeamUseCase interface {
	AbstractGenericUseCase[*model.Team, *model.TeamListOption]

	AddMember(ctx context.Context, teamID, userID string) error
	DeleteMember(ctx context.Context, teamID, userID string) error
}

type teamUseCase struct {
	baseGenericUseCase[*model.Team, *model.TeamListOption]
	svc         service.TeamService
	userService service.UserService
}

func NewTeamUseCase(svc service.TeamService, userService service.UserService) TeamUseCase {
	base := newGenericUseCase(svc)
	return &teamUseCase{baseGenericUseCase: base, svc: svc, userService: userService}
}

var _ TeamUseCase = &teamUseCase{}

func (uc *teamUseCase) AddMember(ctx context.Context, teamID, userID string) error {
	user, err := uc.userService.Get(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}
	return uc.svc.AddMember(ctx, teamID, user)
}
func (uc *teamUseCase) DeleteMember(ctx context.Context, teamID, userID string) error {
	user, err := uc.userService.Get(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}
	return uc.svc.DeleteMember(ctx, teamID, user)
}

package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type TeamService struct {
	repo repository.TeamRepository
	baseGenericService[*model.Team, *model.TeamListOption]
	userService *UserService
}

func NewTeamService(repo repository.TeamRepository, userService UserService) TeamService {
	base := newGenericService(repo)
	return TeamService{repo: repo, baseGenericService: base, userService: &userService}
}

func (s *TeamService) AddMember(ctx context.Context, teamID, userID string) error {
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to get user:%s", userID)
	}
	return s.repo.AddMember(ctx, teamID, user)
}

func (s *TeamService) DeleteMember(ctx context.Context, teamID, userID string) error {
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteMember(ctx, teamID, user)
}

package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type TeamService interface {
	AbstractGenericService[*model.Team, *model.TeamListOption]

	AddMember(ctx context.Context, teamID, userID string) error
	DeleteMember(ctx context.Context, teamID, userID string) error
}

type teamService struct {
	repo repository.TeamRepository
	baseGenericService[*model.Team, *model.TeamListOption]
	userService UserService
}

func NewTeamService(repo repository.TeamRepository, userService UserService) TeamService {
	base := newGenericService(repo)
	return &teamService{repo: repo, baseGenericService: base, userService: userService}
}

func (s *teamService) AddMember(ctx context.Context, teamID, userID string) error {
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to get user:%s", userID)
	}
	return s.repo.AddMember(ctx, teamID, user)
}

func (s *teamService) DeleteMember(ctx context.Context, teamID, userID string) error {
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteMember(ctx, teamID, user)
}

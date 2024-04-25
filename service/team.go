package service

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type TeamService interface {
	AbstractGenericService[*model.Team, *model.TeamListOption]

	AddMember(ctx context.Context, teamID string, user *model.User) error
	DeleteMember(ctx context.Context, teamID string, user *model.User) error
}

type teamService struct {
	repo repository.TeamRepository
	baseGenericService[*model.Team, *model.TeamListOption]
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	base := newGenericService(repo)
	return &teamService{repo: repo, baseGenericService: base}
}

func (s *teamService) AddMember(ctx context.Context, teamID string, user *model.User) error {
	return s.repo.AddMember(ctx, teamID, user)
}

func (s *teamService) DeleteMember(ctx context.Context, teamID string, user *model.User) error {
	return s.repo.DeleteMember(ctx, teamID, user)
}

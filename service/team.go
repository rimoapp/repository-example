package service

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type TeamService struct {
	Repo repository.TeamRepository
	baseGenericService[*model.Team, *model.TeamListOption]
}

func NewTeamService(repo repository.TeamRepository) *TeamService {
	base := newGenericService(repo)
	return &TeamService{Repo: repo, baseGenericService: base}
}

func (s *TeamService) List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error) {
	return s.Repo.List(ctx, opts)
}

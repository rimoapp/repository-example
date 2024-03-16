package service

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type UserService struct {
	Repo repository.UserRepository
	baseGenericService[*model.User, *model.UserListOption]
}

func NewUserService(repo repository.UserRepository) UserService {
	base := newGenericService(repo)
	return UserService{Repo: repo, baseGenericService: base}
}

func (s *UserService) List(ctx context.Context, opts *model.UserListOption) ([]*model.User, error) {
	return s.Repo.List(ctx, opts)
}

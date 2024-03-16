package service

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type UserService interface {
	AbstractGenericService[*model.User, *model.UserListOption]
}

type userService struct {
	repo repository.UserRepository
	baseGenericService[*model.User, *model.UserListOption]
}

func NewUserService(repo repository.UserRepository) UserService {
	base := newGenericService(repo)
	return &userService{repo: repo, baseGenericService: base}
}

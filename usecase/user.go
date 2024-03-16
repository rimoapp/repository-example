package usecase

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
)

type UserUseCase interface {
	AbstractGenericUseCase[*model.User, *model.UserListOption]
}

type userUseCase struct {
	baseGenericUseCase[*model.User, *model.UserListOption]
}

func NewUserUseCase(svc service.UserService) UserUseCase {
	base := newGenericUseCase(svc)
	return &userUseCase{baseGenericUseCase: base}
}

var _ UserUseCase = &userUseCase{}

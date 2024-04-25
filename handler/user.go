package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/usecase"
)

type UserHandler struct {
	baseGenericHandler[*model.User, *model.UserListOption]
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) UserHandler {
	handler := NewGenericHandler(useCase, "userID")
	return UserHandler{baseGenericHandler: handler, useCase: useCase}
}

func (h *UserHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/users", h.List)
	group.GET("/users/:userID", h.Get)
	group.POST("/users", h.Create)
	group.DELETE("/users/:userID", h.Delete)
	group.PATCH("/users/:userID", h.Update)
}

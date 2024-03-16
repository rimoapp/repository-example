package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/usecase"
)

type UserHandler struct {
	baseGenericHandler[*model.User, *model.UserListOption]
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	handler := NewGenericHandler(useCase, "userID")
	return &UserHandler{baseGenericHandler: *handler, useCase: useCase}
}

func (h *UserHandler) List(c *gin.Context) {
	opts := &model.UserListOption{}
	if err := c.ShouldBind(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	opts.UserID = c.GetString("user_id")
	entities, err := h.useCase.List(c, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	filtered := []*model.User{}
	for _, entity := range entities {
		if entity.IsAuthorized(c.GetString("user_id")) {
			filtered = append(filtered, entity)
		}
	}
	c.JSON(http.StatusOK, filtered)
}

func (h *UserHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/users", h.List)
	group.GET("/users/:userID", h.Get)
	group.POST("/users", h.Create)
	group.DELETE("/users/:userID", h.Delete)
	group.PATCH("/users/:userID", h.Update)
}

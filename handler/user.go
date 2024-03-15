package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
)

type UserHandler struct {
	BaseGenericHandler[*model.User]
	Service *service.UserService
}

func NewUserHandler(opts repository.NewRepositoryOption) *UserHandler {
	repo := repository.NewUserRepository(opts)
	return newUserHandler(repo)
}
func newUserHandler(repo repository.UserRepository) *UserHandler {
	svc := service.NewUserService(repo)
	handler := NewGenericHandler(svc, "userID")
	return &UserHandler{BaseGenericHandler: *handler, Service: svc}
}

func (h *UserHandler) List(c *gin.Context) {
	opts := &model.UserListOption{}
	if err := c.ShouldBind(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	opts.UserID = c.GetString("user_id")
	entities, err := h.Service.List(c, opts)
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

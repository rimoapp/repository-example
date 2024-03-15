package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
)

type {{.ModelName}}Handler struct {
	BaseGenericHandler[*model.{{.ModelName}}]
	Service *service.{{.ModelName}}Service
}

func New{{.ModelName}}Handler(opts repository.NewRepositoryOption) *{{.ModelName}}Handler {
	repo := repository.New{{.ModelName}}Repository(opts)
	return new{{.ModelName}}Handler(repo)
}
func new{{.ModelName}}Handler(repo repository.{{.ModelName}}Repository) *{{.ModelName}}Handler {
	svc := service.New{{.ModelName}}Service(repo)
	handler := NewGenericHandler(svc, "{{.Snake}}_id")
	return &{{.ModelName}}Handler{BaseGenericHandler: *handler, Service: svc}
}

func (h *{{.ModelName}}Handler) List(c *gin.Context) {
	opts := &model.{{.ModelName}}ListOption{}
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
	filtered := []*model.{{.ModelName}}{}
	for _, entity := range entities {
		if entity.IsAuthorized(c.GetString("user_id")) {
			filtered = append(filtered, entity)
		}
	}
	c.JSON(http.StatusOK, filtered)
}

func (h *{{.ModelName}}Handler) SetRouter(group *gin.RouterGroup) {
	group.GET("/{{.BasePath}}", h.List)
	group.GET("/{{.BasePath}}/:{{.Snake}}_id", h.Get)
	group.POST("/{{.BasePath}}", h.Create)
	group.DELETE("/{{.BasePath}}/:{{.Snake}}_id", h.Delete)
	group.PATCH("/{{.BasePath}}/:{{.Snake}}_id", h.Update)
}

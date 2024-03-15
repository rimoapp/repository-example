package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
)
type {{.ModelName}}Handler struct {
	BaseGenericHandler[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
	Service *service.{{.ModelName}}Service
}

func New{{.ModelName}}Handler(opts repository.NewRepositoryOption) *{{.ModelName}}Handler {
	repo := repository.New{{.ModelName}}Repository(opts)
	return new{{.ModelName}}Handler(repo)
}
func new{{.ModelName}}Handler(repo repository.{{.ModelName}}Repository) *{{.ModelName}}Handler {
	svc := service.New{{.ModelName}}Service(repo)
	handler := NewGenericHandler(svc, "{{.Snake}}ID")
	return &{{.ModelName}}Handler{BaseGenericHandler: *handler, Service: svc}
}

func (h *{{.ModelName}}Handler) SetRouter(group *gin.RouterGroup) {
	group.GET("/{{.BasePath}}", h.List)
	group.GET("/{{.BasePath}}/:{{.Snake}}ID", h.Get)
	group.POST("/{{.BasePath}}", h.Create)
	group.DELETE("/{{.BasePath}}/:{{.Snake}}ID", h.Delete)
	group.PATCH("/{{.BasePath}}/:{{.Snake}}ID", h.Update)
}

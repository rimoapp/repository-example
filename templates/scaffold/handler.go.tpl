package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
)

type {{.ModelName}}Handler struct {
	baseGenericHandler[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
	svc service.{{.ModelName}}Service
}


func New{{.ModelName}}Handler(svc service.{{.ModelName}}Service) *{{.ModelName}}Handler {
	handler := NewGenericHandler(&svc, "{{.Snake}}ID")
	return &{{.ModelName}}Handler{baseGenericHandler: *handler, svc: svc}
}

func (h *{{.ModelName}}Handler) SetRouter(group *gin.RouterGroup) {
	group.GET("/{{.BasePath}}", h.List)
	group.GET("/{{.BasePath}}/:{{.Snake}}ID", h.Get)
	group.POST("/{{.BasePath}}", h.Create)
	group.DELETE("/{{.BasePath}}/:{{.Snake}}ID", h.Delete)
	group.PATCH("/{{.BasePath}}/:{{.Snake}}ID", h.Update)
}

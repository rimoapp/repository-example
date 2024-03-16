package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/usecase"
)

type {{.ModelName}}Handler struct {
	baseGenericHandler[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
	useCase usecase.{{.ModelName}}UseCase
}


func New{{.ModelName}}Handler(useCase usecase.{{.ModelName}}UseCase) *{{.ModelName}}Handler {
	handler := NewGenericHandler(&svc, "{{.Snake}}ID")
	return &{{.ModelName}}Handler{baseGenericHandler: *handler, useCase: useCase}
}

func (h *{{.ModelName}}Handler) SetRouter(group *gin.RouterGroup) {
	group.GET("/{{.BasePath}}", h.List)
	group.GET("/{{.BasePath}}/:{{.Snake}}ID", h.Get)
	group.POST("/{{.BasePath}}", h.Create)
	group.DELETE("/{{.BasePath}}/:{{.Snake}}ID", h.Delete)
	group.PATCH("/{{.BasePath}}/:{{.Snake}}ID", h.Update)
}

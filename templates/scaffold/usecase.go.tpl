package usecase

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
)

type {{.ModelName}}UseCase interface {
	AbstractGenericUseCase[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
}

type {{.LowerCamelCase}}UseCase struct {
	baseGenericUseCase[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
    svc service.{{.ModelName}}Service
}

func New{{.ModelName}}UseCase(svc service.{{.ModelName}}Service) {{.ModelName}}UseCase {
	base := newGenericUseCase(svc)
	return &{{.LowerCamelCase}}UseCase{baseGenericUseCase: base, svc: svc}
}

var _ {{.ModelName}}UseCase = &{{.LowerCamelCase}}UseCase{}

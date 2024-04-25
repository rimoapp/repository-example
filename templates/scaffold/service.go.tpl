package service

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type {{.ModelName}}Service interface {
	AbstractGenericService[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
}

type {{.LowerCamelCase}}Service struct {
	repo repository.{{.ModelName}}Repository
	baseGenericService[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
}

func New{{.ModelName}}Service(repo repository.{{.ModelName}}Repository) {{.ModelName}}Service {
	base := newGenericService(repo)
	return &{{.LowerCamelCase}}Service{repo: repo, baseGenericService: base}
}

var _ {{.ModelName}}Service = &{{.LowerCamelCase}}Service{}

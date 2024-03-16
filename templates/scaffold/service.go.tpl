package service

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type {{.ModelName}}Service struct {
	Repo repository.{{.ModelName}}Repository
	baseGenericService[*model.{{.ModelName}}, *model.{{.ModelName}}ListOption]
}

func New{{.ModelName}}Service(repo repository.{{.ModelName}}Repository) {{.ModelName}}Service {
	base := newGenericService(repo)
	return {{.ModelName}}Service{Repo: repo, baseGenericService: base}
}

package service

import (
	"context"
	"time"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type {{.ModelName}}Service struct {
	Repo repository.{{.ModelName}}Repository
	BaseGenericService[*model.{{.ModelName}}]
}

func New{{.ModelName}}Service(repo repository.{{.ModelName}}Repository) *{{.ModelName}}Service {
	return &{{.ModelName}}Service{Repo: repo}
}

func (s *{{.ModelName}}Service) List(ctx context.Context, opts *model.{{.ModelName}}ListOption) ([]*model.{{.ModelName}}, error) {
	return s.Repo.List(ctx, opts)
}

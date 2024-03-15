package service

import (
	"context"
	"time"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type {{.ModelName}}Service struct {
	Repo repository.{{.ModelName}}Repository
}


func New{{.ModelName}}Service(repo repository.{{.ModelName}}Repository) *{{.ModelName}}Service {
	return &{{.ModelName}}Service{Repo: repo}
}

func (s *{{.ModelName}}Service) Get(ctx context.Context, id string, opts model.AbstractDataEntityOption) (*model.{{.ModelName}}, error) {
	return s.Repo.Get(ctx, id, opts)
}

func (s *{{.ModelName}}Service) Delete(ctx context.Context, id string, opts model.AbstractDataEntityOption) error {
	return s.Repo.Delete(ctx, id, opts)
}

func (s *{{.ModelName}}Service) Create(ctx context.Context, object *model.{{.ModelName}}, opts model.AbstractDataEntityOption) (string, error) {
	now := time.Now()
	object.CreatedAt = now
	object.UpdatedAt = now
	return s.Repo.Create(ctx, object, opts)
}

func (s *{{.ModelName}}Service) Update(ctx context.Context, id string, keyValues map[string]interface{}, opts model.AbstractDataEntityOption) error {
	return s.Repo.Update(ctx, id, keyValues, opts)
}

func (s *{{.ModelName}}Service) List(ctx context.Context, opts model.AbstractListOption[*model.{{.ModelName}}]) ([]*model.{{.ModelName}}, error) {
	return s.Repo.List(ctx, opts)
}

func (s *{{.ModelName}}Service) Set(ctx context.Context, id string, object *model.{{.ModelName}}, opts model.AbstractDataEntityOption) error {
	return s.Repo.Set(ctx, id, object, opts)
}



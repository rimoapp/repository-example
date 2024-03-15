package service

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type AbstractGenericService[T model.AbstractDataEntity] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object T) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	Set(ctx context.Context, id string, object T) error
}

// BaseGenericService is a service for generic repository
type BaseGenericService[T model.AbstractDataEntity] struct {
	Repo repository.AbstractGenericRepository[T]
}

var _ AbstractGenericService[model.AbstractDataEntity] = &BaseGenericService[model.AbstractDataEntity]{}

func NewGenericService[T model.AbstractDataEntity](repo repository.AbstractGenericRepository[T]) AbstractGenericService[T] {
	return &BaseGenericService[T]{Repo: repo}
}

func (s *BaseGenericService[T]) Get(ctx context.Context, id string) (T, error) {
	return s.Repo.Get(ctx, id)
}

func (s *BaseGenericService[T]) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *BaseGenericService[T]) Create(ctx context.Context, object T) (string, error) {

	return s.Repo.Create(ctx, object)
}

func (s *BaseGenericService[T]) Update(ctx context.Context, id string, keyValues map[string]interface{}) error {
	return s.Repo.Update(ctx, id, keyValues)
}

func (s *BaseGenericService[T]) Set(ctx context.Context, id string, object T) error {
	return s.Repo.Set(ctx, id, object)
}

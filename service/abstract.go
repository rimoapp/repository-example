package service

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type AbstractGenericService[T model.AbstractDataEntity, U model.AbstractListOption] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object T) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	Set(ctx context.Context, id string, object T) error
	List(ctx context.Context, opts U) ([]T, error)
}

// baseGenericService is a service for generic repository
type baseGenericService[T model.AbstractDataEntity, U model.AbstractListOption] struct {
	Repo repository.AbstractGenericRepository[T, U]
}

var _ AbstractGenericService[model.AbstractDataEntity, model.AbstractListOption] = &baseGenericService[model.AbstractDataEntity, model.AbstractListOption]{}

func newGenericService[T model.AbstractDataEntity, U model.AbstractListOption](repo repository.AbstractGenericRepository[T, U]) baseGenericService[T, U] {
	return baseGenericService[T, U]{Repo: repo}
}

func (s *baseGenericService[T, U]) Get(ctx context.Context, id string) (T, error) {
	return s.Repo.Get(ctx, id)
}

func (s *baseGenericService[T, U]) Delete(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *baseGenericService[T, U]) Create(ctx context.Context, object T) (string, error) {

	return s.Repo.Create(ctx, object)
}

func (s *baseGenericService[T, U]) Update(ctx context.Context, id string, keyValues map[string]interface{}) error {
	return s.Repo.Update(ctx, id, keyValues)
}

func (s *baseGenericService[T, U]) Set(ctx context.Context, id string, object T) error {
	return s.Repo.Set(ctx, id, object)
}

func (s *baseGenericService[T, U]) List(ctx context.Context, opts U) ([]T, error) {
	return s.Repo.List(ctx, opts)
}

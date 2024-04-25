package service

import (
	"context"
	"time"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type AbstractGenericService[T model.AbstractEntity, U model.AbstractListOption] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object T) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	Set(ctx context.Context, id string, object T) error
	List(ctx context.Context, opts U) ([]T, error)
}

// baseGenericService is a service for generic repository
type baseGenericService[T model.AbstractEntity, U model.AbstractListOption] struct {
	repo repository.AbstractGenericRepository[T, U]
}

var _ AbstractGenericService[model.AbstractEntity, model.AbstractListOption] = &baseGenericService[model.AbstractEntity, model.AbstractListOption]{}

func newGenericService[T model.AbstractEntity, U model.AbstractListOption](repo repository.AbstractGenericRepository[T, U]) baseGenericService[T, U] {
	return baseGenericService[T, U]{repo: repo}
}

func (s *baseGenericService[T, U]) Get(ctx context.Context, id string) (T, error) {
	return s.repo.Get(ctx, id)
}

func (s *baseGenericService[T, U]) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *baseGenericService[T, U]) Create(ctx context.Context, object T) (string, error) {
	object.HooksBeforeCreate(time.Now())
	return s.repo.Create(ctx, object)
}

func (s *baseGenericService[T, U]) Update(ctx context.Context, id string, keyValues map[string]interface{}) error {
	return s.repo.Update(ctx, id, keyValues)
}

func (s *baseGenericService[T, U]) Set(ctx context.Context, id string, object T) error {
	return s.repo.Set(ctx, id, object)
}

func (s *baseGenericService[T, U]) List(ctx context.Context, opts U) ([]T, error) {
	return s.repo.List(ctx, opts)
}

package usecase

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
)

type AbstractGenericUseCase[T model.AbstractEntity, U model.AbstractListOption] interface {
	Get(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object T) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	Set(ctx context.Context, id string, object T) error
	List(ctx context.Context, opts U) ([]T, error)
}

// baseGenericUseCase is a service for generic repository
type baseGenericUseCase[T model.AbstractEntity, U model.AbstractListOption] struct {
	svc service.AbstractGenericService[T, U]
}

var _ AbstractGenericUseCase[model.AbstractEntity, model.AbstractListOption] = &baseGenericUseCase[model.AbstractEntity, model.AbstractListOption]{}

func newGenericUseCase[T model.AbstractEntity, U model.AbstractListOption](svc service.AbstractGenericService[T, U]) baseGenericUseCase[T, U] {
	return baseGenericUseCase[T, U]{svc: svc}
}

func (s *baseGenericUseCase[T, U]) Get(ctx context.Context, id string) (T, error) {
	return s.svc.Get(ctx, id)
}

func (s *baseGenericUseCase[T, U]) Delete(ctx context.Context, id string) error {
	return s.svc.Delete(ctx, id)
}

func (s *baseGenericUseCase[T, U]) Create(ctx context.Context, object T) (string, error) {
	return s.svc.Create(ctx, object)
}

func (s *baseGenericUseCase[T, U]) Update(ctx context.Context, id string, keyValues map[string]interface{}) error {
	return s.svc.Update(ctx, id, keyValues)
}

func (s *baseGenericUseCase[T, U]) Set(ctx context.Context, id string, object T) error {
	return s.svc.Set(ctx, id, object)
}

func (s *baseGenericUseCase[T, U]) List(ctx context.Context, opts U) ([]T, error) {
	return s.svc.List(ctx, opts)
}

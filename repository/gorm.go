package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/gogo/status"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type gormGenericRepository[T model.AbstractEntity] struct {
	client *gorm.DB
}

func newGormGenericRepository[T model.AbstractEntity](client *gorm.DB) *gormGenericRepository[T] {
	var object T
	client.AutoMigrate(object)
	return &gormGenericRepository[T]{
		client: client,
	}
}

func (r *gormGenericRepository[T]) Get(ctx context.Context, id string) (T, error) {
	entity := createNewInstance[T]()
	uid, err := strconv.Atoi(id)
	if err != nil {
		return entity, errors.Wrap(err, "failed to convert id to int")
	}
	result := r.client.First(entity, uid)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity, status.New(codes.NotFound, "not found").Err()
		}
		return entity, result.Error
	}
	entity.SetID(entity.GetID())
	return entity, nil
}

func (r *gormGenericRepository[T]) Delete(ctx context.Context, id string) error {
	entity := createNewInstance[T]()
	entity.SetID(id)
	result := r.client.Delete(entity)
	return errors.Wrap(result.Error, "failed to delete")
}

func (r *gormGenericRepository[T]) Create(ctx context.Context, object T) (string, error) {
	object.BeforeCreate(time.Now())
	result := r.client.Create(object)
	return object.GetID(), result.Error
}

func (r *gormGenericRepository[T]) Update(ctx context.Context, id string, keyValues map[string]interface{}) error {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return errors.Wrap(err, "failed to convert id to int")
	}
	var object T
	result := r.client.Model(object).Where("uid = ?", uid).Updates(keyValues)
	return result.Error
}

func (r *gormGenericRepository[T]) list(query *gorm.DB) ([]T, error) {
	objects := []T{}
	result := query.Find(&objects)
	for _, entity := range objects {
		entity.SetID(entity.GetID())
	}
	return objects, result.Error
}

func (r *gormGenericRepository[T]) Set(ctx context.Context, id string, object T) error {
	object.SetID(id)
	result := r.client.Save(object)
	return result.Error
}

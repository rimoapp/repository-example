package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
)

// TODO: Delete unused functions
type UserRepository interface {
	Get(ctx context.Context, id string) (*model.User, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.User) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.UserListOption) ([]*model.User, error)
	Set(ctx context.Context, id string, object *model.User) error
}

const usersCollectionPath = "users"

func NewUserRepository(opts NewRepositoryOption) UserRepository {
	if opts.DBClient != nil {
		return &gormUserRepository{
			gormGenericRepository: *newGormGenericRepository[*model.User](opts.DBClient),
		}
	}
	return &firestoreUserRepository{
		firestoreGenericRepository: *newFirestoreGenericRepository[*model.User](opts.FirestoreClient, usersCollectionPath),
	}
}

type firestoreUserRepository struct {
	firestoreGenericRepository[*model.User]
}

func (r *firestoreUserRepository) List(ctx context.Context, opts *model.UserListOption) ([]*model.User, error) {
	collectionPath := r.collectionPath
	if collectionPath == "" {
		return nil, errors.New("collection path is empty")
	}
	query := r.client.Collection(collectionPath).Query
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id", "==", opts.UserID)
	}
	return r.firestoreGenericRepository.list(ctx, query)
}

type gormUserRepository struct {
	gormGenericRepository[*model.User]
}

func (r *gormUserRepository) List(ctx context.Context, opts *model.UserListOption) ([]*model.User, error) {
	query := r.client
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id = ?", opts.UserID)
	}
	return r.gormGenericRepository.list(query)
}

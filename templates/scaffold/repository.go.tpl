package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
)

type {{.ModelName}}Repository interface {
	Get(ctx context.Context, id string) (*model.{{.ModelName}}, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.{{.ModelName}}) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.{{.ModelName}}ListOption) ([]*model.{{.ModelName}}, error)
	Set(ctx context.Context, id string, object *model.{{.ModelName}}) error
}

const {{.TableName}}CollectionPath = "{{.TableName}}"

func New{{.ModelName}}Repository(opts *NewRepositoryOption) {{.ModelName}}Repository {
	if opts.dbClient != nil {
		return &gorm{{.ModelName}}Repository{
			gormGenericRepository: *newGormGenericRepository[*model.{{.ModelName}}](opts.dbClient),
		}
	}
	return &firestore{{.ModelName}}Repository{
		firestoreGenericRepository: *newFirestoreGenericRepository[*model.{{.ModelName}}](opts.firestoreClient, {{.TableName}}CollectionPath),
	}
}

type firestore{{.ModelName}}Repository struct {
	firestoreGenericRepository[*model.{{.ModelName}}]
}

func (r *firestore{{.ModelName}}Repository) List(ctx context.Context, opts *model.{{.ModelName}}ListOption) ([]*model.{{.ModelName}}, error) {
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

type gorm{{.ModelName}}Repository struct {
	gormGenericRepository[*model.{{.ModelName}}]
}

func (r *gorm{{.ModelName}}Repository) List(ctx context.Context, opts *model.{{.ModelName}}ListOption) ([]*model.{{.ModelName}}, error) {
	query := r.client
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id = ?", opts.UserID)
	}
	return r.gormGenericRepository.list(query)
}

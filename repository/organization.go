package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
)

type OrganizationRepository interface {
	Get(ctx context.Context, id string) (*model.Organization, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.Organization) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.OrganizationListOption) ([]*model.Organization, error)
	Set(ctx context.Context, id string, object *model.Organization) error
}

const organizationsCollectionPath = "organizations"

func NewOrganizationRepository(opts *NewRepositoryOption) OrganizationRepository {
	if opts.dbClient != nil {
		return &gormOrganizationRepository{
			gormGenericRepository: *newGormGenericRepository[*model.Organization](opts.dbClient),
		}
	}
	return &firestoreOrganizationRepository{
		firestoreGenericRepository: *newFirestoreGenericRepository[*model.Organization](opts.firestoreClient, organizationsCollectionPath),
	}
}

type firestoreOrganizationRepository struct {
	firestoreGenericRepository[*model.Organization]
}

func (r *firestoreOrganizationRepository) List(ctx context.Context, opts *model.OrganizationListOption) ([]*model.Organization, error) {
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

type gormOrganizationRepository struct {
	gormGenericRepository[*model.Organization]
}

func (r *gormOrganizationRepository) List(ctx context.Context, opts *model.OrganizationListOption) ([]*model.Organization, error) {
	query := r.client
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id = ?", opts.UserID)
	}
	return r.gormGenericRepository.list(query)
}

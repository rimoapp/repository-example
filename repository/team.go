package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
)

// TODO: Delete unused functions
type TeamRepository interface {
	Get(ctx context.Context, id string) (*model.Team, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.Team) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error)
	Set(ctx context.Context, id string, object *model.Team) error
}

const teamsCollectionPath = "teams"

func NewTeamRepository(opts NewRepositoryOption) TeamRepository {
	if opts.DBClient != nil {
		return &gormTeamRepository{
			gormGenericRepository: *newGormGenericRepository[*model.Team](opts.DBClient),
		}
	}
	return &firestoreTeamRepository{
		firestoreGenericRepository: *newFirestoreGenericRepository[*model.Team](opts.FirestoreClient, teamsCollectionPath),
	}
}

type firestoreTeamRepository struct {
	firestoreGenericRepository[*model.Team]
}

func (r *firestoreTeamRepository) List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error) {
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

type gormTeamRepository struct {
	gormGenericRepository[*model.Team]
}

func (r *gormTeamRepository) List(ctx context.Context, opts *model.TeamListOption) ([]*model.Team, error) {
	query := r.client
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id = ?", opts.UserID)
	}
	return r.gormGenericRepository.list(query)
}

package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
)

type NoteRepository interface {
	Get(ctx context.Context, id string) (*model.Note, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, object *model.Note) (string, error)
	Update(ctx context.Context, id string, keyValues map[string]interface{}) error
	List(ctx context.Context, opts *model.NoteListOption) ([]*model.Note, error)
	Set(ctx context.Context, id string, object *model.Note) error
}

const notesCollectionPath = "notes"

func NewNoteRepository(opts *NewRepositoryOption) NoteRepository {
	if opts.dbClient != nil {
		return &gormNoteRepository{
			gormGenericRepository: *newGormGenericRepository[*model.Note](opts.dbClient),
		}
	}
	return &firestoreNoteRepository{
		firestoreGenericRepository: *newFirestoreGenericRepository[*model.Note](opts.firestoreClient, notesCollectionPath),
	}
}

type firestoreNoteRepository struct {
	firestoreGenericRepository[*model.Note]
}

func (r *firestoreNoteRepository) List(ctx context.Context, opts *model.NoteListOption) ([]*model.Note, error) {
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

type gormNoteRepository struct {
	gormGenericRepository[*model.Note]
}

func (r *gormNoteRepository) List(ctx context.Context, opts *model.NoteListOption) ([]*model.Note, error) {
	query := r.client
	// NOTE: implement where clause
	if opts.UserID != "" {
		query = query.Where("user_id = ?", opts.UserID)
	}
	return r.gormGenericRepository.list(query)
}

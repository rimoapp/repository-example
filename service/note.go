package service

import (
	"context"

	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type NoteService struct {
	Repo repository.NoteRepository
	BaseGenericService[*model.Note]
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	base := NewGenericService(repo)
	return &NoteService{Repo: repo, BaseGenericService: base}
}

func (s *NoteService) List(ctx context.Context, opts *model.NoteListOption) ([]*model.Note, error) {
	return s.Repo.List(ctx, opts)
}

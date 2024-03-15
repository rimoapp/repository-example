package service

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type NoteService struct {
	Repo repository.NoteRepository
	BaseGenericService[*model.Note, *model.NoteListOption]
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	base := NewGenericService(repo)
	return &NoteService{Repo: repo, BaseGenericService: base}
}

package service

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type NoteService struct {
	Repo repository.NoteRepository
	baseGenericService[*model.Note, *model.NoteListOption]
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	base := newGenericService(repo)
	return NoteService{Repo: repo, baseGenericService: base}
}

var _ AbstractGenericService[*model.Note, *model.NoteListOption] = &NoteService{}

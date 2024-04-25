package service

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
)

type NoteService interface {
	AbstractGenericService[*model.Note, *model.NoteListOption]
}

type noteService struct {
	repo repository.NoteRepository
	baseGenericService[*model.Note, *model.NoteListOption]
}

func NewNoteService(repo repository.NoteRepository) NoteService {
	base := newGenericService(repo)
	return &noteService{repo: repo, baseGenericService: base}
}

var _ NoteService = &noteService{}

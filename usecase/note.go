package usecase

import (
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
)

type NoteUseCase interface {
	AbstractGenericUseCase[*model.Note, *model.NoteListOption]
}

type noteUseCase struct {
	baseGenericUseCase[*model.Note, *model.NoteListOption]
}

func NewNoteUseCase(svc service.NoteService) NoteUseCase {
	base := newGenericUseCase(svc)
	return &noteUseCase{baseGenericUseCase: base}
}

var _ NoteUseCase = &noteUseCase{}

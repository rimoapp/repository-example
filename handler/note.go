package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/usecase"
)

type NoteHandler struct {
	baseGenericHandler[*model.Note, *model.NoteListOption]
	useCase usecase.NoteUseCase
}

func NewNoteHandler(useCase usecase.NoteUseCase) *NoteHandler {
	handler := NewGenericHandler(useCase, "noteID")
	return &NoteHandler{baseGenericHandler: *handler, useCase: useCase}
}

func (h *NoteHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/notes", h.List)
	group.GET("/notes/:noteID", h.Get)
	group.POST("/notes", h.Create)
	group.DELETE("/notes/:noteID", h.Delete)
	group.PATCH("/notes/:noteID", h.Update)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
)

type NoteHandler struct {
	baseGenericHandler[*model.Note, *model.NoteListOption]
	svc *service.NoteService
}

func NewNoteHandler(svc service.NoteService) *NoteHandler {
	handler := NewGenericHandler(&svc, "noteID")
	return &NoteHandler{baseGenericHandler: *handler, svc: &svc}
}

func (h *NoteHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/notes", h.List)
	group.GET("/notes/:noteID", h.Get)
	group.POST("/notes", h.Create)
	group.DELETE("/notes/:noteID", h.Delete)
	group.PATCH("/notes/:noteID", h.Update)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
)
type NoteHandler struct {
	BaseGenericHandler[*model.Note, *model.NoteListOption]
	Service *service.NoteService
}

func NewNoteHandler(opts repository.NewRepositoryOption) *NoteHandler {
	repo := repository.NewNoteRepository(opts)
	return newNoteHandler(repo)
}
func newNoteHandler(repo repository.NoteRepository) *NoteHandler {
	svc := service.NewNoteService(repo)
	handler := NewGenericHandler(svc, "noteID")
	return &NoteHandler{BaseGenericHandler: *handler, Service: svc}
}

func (h *NoteHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/notes", h.List)
	group.GET("/notes/:noteID", h.Get)
	group.POST("/notes", h.Create)
	group.DELETE("/notes/:noteID", h.Delete)
	group.PATCH("/notes/:noteID", h.Update)
}

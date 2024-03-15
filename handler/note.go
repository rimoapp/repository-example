package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
)

type NoteHandler struct {
	BaseGenericHandler[*model.Note]
	Service *service.NoteService
}

func NewNoteHandler(opts repository.NewRepositoryOption) *NoteHandler {
	repo := repository.NewNoteRepository(opts)
	return newNoteHandler(repo)
}
func newNoteHandler(repo repository.NoteRepository) *NoteHandler {
	svc := service.NewNoteService(repo)
	handler := NewGenericHandler(svc, "note_id")
	return &NoteHandler{BaseGenericHandler: *handler, Service: svc}
}

func (h *NoteHandler) List(c *gin.Context) {
	opts := &model.NoteListOption{}
	if err := c.ShouldBind(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	opts.UserID = c.GetString("user_id")
	entities, err := h.Service.List(c, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	filtered := []*model.Note{}
	for _, entity := range entities {
		if entity.IsAuthorized(c.GetString("user_id")) {
			filtered = append(filtered, entity)
		}
	}
	c.JSON(http.StatusOK, filtered)
}

func (h *NoteHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/notes", h.List)
	group.GET("/notes/:note_id", h.Get)
	group.POST("/notes", h.Create)
	group.DELETE("/notes/:note_id", h.Delete)
	group.PATCH("/notes/:note_id", h.Update)
}

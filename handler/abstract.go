package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type baseGenericHandler[T model.AbstractAssociatedEntity, U model.AbstractListOption] struct {
	svc     service.AbstractGenericService[T, U]
	idParam string
}

func NewGenericHandler[T model.AbstractAssociatedEntity, U model.AbstractListOption](service service.AbstractGenericService[T, U], idParam string) *baseGenericHandler[T, U] {
	return &baseGenericHandler[T, U]{svc: service, idParam: idParam}
}

func (h *baseGenericHandler[T, U]) authWithEntity(c *gin.Context) (T, bool) {
	id := c.Param(h.idParam)
	entity, err := h.svc.Get(c, id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return entity, false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to get").Error()})
		return entity, false
	}
	if !entity.IsAuthorized(c.GetString("user_id")) {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return entity, false
	}
	return entity, true
}

func (h *baseGenericHandler[T, U]) Get(c *gin.Context) {
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *baseGenericHandler[T, U]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBind(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	entity.SetCreatorID(c.GetString("user_id"))
	id, err := h.svc.Create(c, entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	entity.SetID(id)
	c.JSON(http.StatusOK, entity)
}

func (h *baseGenericHandler[T, U]) Delete(c *gin.Context) {
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	if err := h.svc.Delete(c, entity.GetID()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *baseGenericHandler[T, U]) Update(c *gin.Context) {
	req := map[string]interface{}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	if err := h.svc.Update(c, entity.GetID(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	entity, err := h.svc.Get(c, entity.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *baseGenericHandler[T, U]) List(c *gin.Context) {
	var opts U
	if err := c.ShouldBind(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	opts.SetUserID(c.GetString("user_id"))
	entities, err := h.svc.List(c, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	filtered := []T{}
	for _, entity := range entities {
		if entity.IsAuthorized(c.GetString("user_id")) {
			filtered = append(filtered, entity)
		}
	}
	c.JSON(http.StatusOK, filtered)
}

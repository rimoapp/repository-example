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

type BaseGenericHandler[T model.AbstractAssociatedEntity] struct {
	service service.AbstractGenericService[T]
	key     string
}

func NewGenericHandler[T model.AbstractAssociatedEntity](service service.AbstractGenericService[T], key string) *BaseGenericHandler[T] {
	return &BaseGenericHandler[T]{service: service, key: key}
}

func (h *BaseGenericHandler[T]) authWithEntity(c *gin.Context) (T, bool) {
	id := c.Param(h.key)
	entity, err := h.service.Get(c, id)
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

func (h *BaseGenericHandler[T]) Get(c *gin.Context) {
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *BaseGenericHandler[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBind(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	entity.SetCreatorID(c.GetString("user_id"))
	id, err := h.service.Create(c, entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	entity.SetID(id)
	c.JSON(http.StatusOK, entity)
}

func (h *BaseGenericHandler[T]) Delete(c *gin.Context) {
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	if err := h.service.Delete(c, entity.GetID()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *BaseGenericHandler[T]) Update(c *gin.Context) {
	req := map[string]interface{}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	if err := h.service.Update(c, entity.GetID(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	entity, err := h.service.Get(c, entity.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	c.JSON(http.StatusOK, entity)
}

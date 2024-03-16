package handler

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type baseGenericHandler[T model.AbstractAssociatedEntity, U model.AbstractListOption] struct {
	useCase usecase.AbstractGenericUseCase[T, U]
	idParam string
}

func NewGenericHandler[T model.AbstractAssociatedEntity, U model.AbstractListOption](useCase usecase.AbstractGenericUseCase[T, U], idParam string) baseGenericHandler[T, U] {
	return baseGenericHandler[T, U]{useCase: useCase, idParam: idParam}
}

func (h *baseGenericHandler[T, U]) authWithEntity(c *gin.Context) (T, bool) {
	id := c.Param(h.idParam)
	entity, err := h.useCase.Get(c, id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return entity, false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to get in auth with entity").Error()})
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
	entity := createNewInstance[T]()
	if err := c.ShouldBind(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}

	entity.SetCreatorID(c.GetString("user_id"))
	id, err := h.useCase.Create(c, entity)
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
	if err := h.useCase.Delete(c, entity.GetID()); err != nil {
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
	if err := h.useCase.Update(c, entity.GetID(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	entity, err := h.useCase.Get(c, entity.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *baseGenericHandler[T, U]) List(c *gin.Context) {
	opts := createNewInstance[U]()
	if err := c.ShouldBind(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}

	opts.SetUserID(c.GetString("user_id"))
	entities, err := h.useCase.List(c, opts)
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

func createNewInstance[T any]() T {
	var entity T
	// reflect.TypeOf を使用して T の型情報を取得
	t := reflect.TypeOf(entity)
	// T がポインタ型の場合、新しいインスタンスを生成
	if t.Kind() == reflect.Ptr {
		// 新しいインスタンスを生成して返す
		newInstance := reflect.New(t.Elem()).Interface()
		return newInstance.(T)
	}
	// ポインタでない場合はそのままデフォルト値を返す
	return entity
}

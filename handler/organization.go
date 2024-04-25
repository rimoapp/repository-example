package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationHandler struct {
	baseGenericHandler[*model.Organization, *model.OrganizationListOption]
	useCase usecase.OrganizationUseCase
}

func NewOrganizationHandler(useCase usecase.OrganizationUseCase) OrganizationHandler {
	handler := NewGenericHandler(useCase, "organizationID")
	return OrganizationHandler{baseGenericHandler: handler, useCase: useCase}
}

func (h *OrganizationHandler) Get(c *gin.Context) {
	req := &model.GetOrganizationOption{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}

	id := c.Param(h.idParam)
	entity, err := h.useCase.GetWithOption(c, id, req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to get").Error()})
		return
	}
	if !entity.IsAuthorized(c.GetString("user_id")) {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *OrganizationHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/organizations", h.List)
	group.GET("/organizations/:organizationID", h.Get)
	group.POST("/organizations", h.Create)
	group.DELETE("/organizations/:organizationID", h.Delete)
	group.PATCH("/organizations/:organizationID", h.Update)
}

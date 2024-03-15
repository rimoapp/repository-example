package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
)

type OrganizationHandler struct {
	BaseGenericHandler[*model.Organization, *model.OrganizationListOption]
	Service *service.OrganizationService
}

func NewOrganizationHandler(opts repository.NewRepositoryOption) *OrganizationHandler {
	repo := repository.NewOrganizationRepository(opts)
	teamsRepo := repository.NewTeamRepository(opts)
	return newOrganizationHandler(repo, teamsRepo)
}
func newOrganizationHandler(repo repository.OrganizationRepository, teamsRepo repository.TeamRepository) *OrganizationHandler {
	teamsSvc := service.NewTeamService(teamsRepo)
	svc := service.NewOrganizationService(repo, teamsSvc)
	handler := NewGenericHandler(svc, "organizationID")
	return &OrganizationHandler{BaseGenericHandler: *handler, Service: svc}
}

func (h *OrganizationHandler) List(c *gin.Context) {
	opts := &model.OrganizationListOption{}
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
	filtered := []*model.Organization{}
	for _, entity := range entities {
		if entity.IsAuthorized(c.GetString("user_id")) {
			filtered = append(filtered, entity)
		}
	}
	c.JSON(http.StatusOK, filtered)
}

func (h *OrganizationHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/organizations", h.List)
	group.GET("/organizations/:organizationID", h.Get)
	group.POST("/organizations", h.Create)
	group.DELETE("/organizations/:organizationID", h.Delete)
	group.PATCH("/organizations/:organizationID", h.Update)
}

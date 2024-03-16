package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rimoapp/repository-example/model"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TeamHandler struct {
	baseGenericHandler[*model.Team, *model.TeamListOption]
	Service *service.TeamService
}

func NewTeamHandler(opts repository.NewRepositoryOption) *TeamHandler {
	repo := repository.NewTeamRepository(opts)
	userRepo := repository.NewUserRepository(opts)
	userSvc := service.NewUserService(userRepo)
	svc := service.NewTeamService(repo, userSvc)
	handler := NewGenericHandler(svc, "teamID")
	return &TeamHandler{baseGenericHandler: *handler, Service: svc}
}

func (h *TeamHandler) List(c *gin.Context) {
	opts := &model.TeamListOption{}
	if err := c.ShouldBind(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	opts.OrganizationID = c.Param("organizationID")
	opts.UserID = c.GetString("user_id")
	entities, err := h.Service.List(c, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	filtered := []*model.Team{}
	for _, entity := range entities {
		if entity.IsAuthorized(c.GetString("user_id")) {
			filtered = append(filtered, entity)
		}
	}
	c.JSON(http.StatusOK, filtered)
}

func (h *TeamHandler) SetRouter(group *gin.RouterGroup) {
	group.GET("/teams", h.List)
	group.GET("/teams/:teamID", h.Get)
	group.POST("/teams", h.Create)
	group.DELETE("/teams/:teamID", h.Delete)
	group.PATCH("/teams/:teamID", h.Update)
	group.POST("/teams/:teamID/members/:userID", h.AddTeamMember)
	group.DELETE("/teams/:teamID/members/:userID", h.DeleteTeamMember)
}

//

func (h *TeamHandler) authWithEntity(c *gin.Context) (*model.Team, bool) {
	id := c.Param("teamID")
	entity, err := h.Service.Get(c, id)
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
	orgID := c.Param("organizationID")
	if orgID != "" && orgID != entity.OrganizationID {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	}
	return entity, true
}

func (h *TeamHandler) Get(c *gin.Context) {
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *TeamHandler) Create(c *gin.Context) {
	var entity *model.Team
	if err := c.ShouldBind(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	entity.OrganizationID = c.Param("organizationID")
	entity.SetCreatorID(c.GetString("user_id"))
	id, err := h.Service.Create(c, entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	entity.SetID(id)
	c.JSON(http.StatusOK, entity)
}

func (h *TeamHandler) Delete(c *gin.Context) {
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	if err := h.Service.Delete(c, entity.GetID()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TeamHandler) Update(c *gin.Context) {
	req := map[string]interface{}{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.Wrap(err, "failed to bind params").Error()})
		return
	}
	entity, authorized := h.authWithEntity(c)
	if !authorized {
		return
	}
	if err := h.Service.Update(c, entity.GetID(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	entity, err := h.Service.Get(c, entity.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to create").Error()})
		return
	}
	c.JSON(http.StatusOK, entity)
}

// AddTeamMember adds a member to the team
func (h *TeamHandler) AddTeamMember(c *gin.Context) {
	teamID := c.Param("teamID")
	userID := c.Param("userID")
	err := h.Service.AddMember(c, teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to add member").Error()})
		return
	}
	c.Status(http.StatusOK)
}

// DeleteTeamMember deletes a member from the team
func (h *TeamHandler) DeleteTeamMember(c *gin.Context) {
	teamID := c.Param("teamID")
	userID := c.Param("userID")
	err := h.Service.DeleteMember(c, teamID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errors.Wrap(err, "failed to delete member").Error()})
		return
	}
	c.Status(http.StatusOK)
}

package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/handler"
	"github.com/rimoapp/repository-example/middleware"
	"github.com/rimoapp/repository-example/repository"
	"github.com/rimoapp/repository-example/service"
	"github.com/rimoapp/repository-example/usecase"
)

type RimoRouter struct {
	http.Handler
}

func NewRouter() (RimoRouter, error) {
	router := gin.New()
	router.Use(gin.Recovery())
	ctx := context.Background()
	repoOpts, err := repository.BuildNewRepositoryOptions(ctx)
	if err != nil {
		return RimoRouter{}, err
	}
	rootGroup := router.Group("/")
	rootGroup.Use(gin.Logger())
	rootGroup.Use(middleware.AuthMock())

	// build repositories
	orgRepository := repository.NewOrganizationRepository(repoOpts)
	teamRepository := repository.NewTeamRepository(repoOpts)
	userRepository := repository.NewUserRepository(repoOpts)
	noteRepository := repository.NewNoteRepository(repoOpts)

	// build services
	noteService := service.NewNoteService(noteRepository)
	userService := service.NewUserService(userRepository)
	teamService := service.NewTeamService(teamRepository)
	orgService := service.NewOrganizationService(orgRepository)

	// build use cases
	organizationUseCase := usecase.NewOrganizationUseCase(orgService, teamService)
	noteUseCase := usecase.NewNoteUseCase(noteService)
	userUseCase := usecase.NewUserUseCase(userService)
	teamUseCase := usecase.NewTeamUseCase(teamService, userService)

	// build handlers
	orgHandler := handler.NewOrganizationHandler(organizationUseCase)
	teamHandler := handler.NewTeamHandler(teamUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	noteHandler := handler.NewNoteHandler(noteUseCase)

	// set routing
	userHandler.SetRouter(rootGroup)
	noteHandler.SetRouter(rootGroup)
	orgHandler.SetRouter(rootGroup)
	orgGroup := rootGroup.Group("/organizations/:organizationID")
	teamHandler.SetRouter(orgGroup)

	rimoRouter := RimoRouter{Handler: router}
	return rimoRouter, nil
}

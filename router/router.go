package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rimoapp/repository-example/handler"
	"github.com/rimoapp/repository-example/middleware"
	"github.com/rimoapp/repository-example/repository"
)

type RimoRouter struct {
	http.Handler
}

func NewRouter() (RimoRouter, error) {
	router := gin.New()
	router.Use(gin.Recovery())
	ctx := context.Background()
	opts, err := repository.BuildNewRepositoryOptions(ctx)
	if err != nil {
		return RimoRouter{}, err
	}
	rootGroup := router.Group("/")
	rootGroup.Use(gin.Logger())
	rootGroup.Use(middleware.AuthMock())

	orgHandler := handler.NewOrganizationHandler(*opts)

	orgHandler.SetRouter(rootGroup)

	rimoRouter := RimoRouter{Handler: router}
	return rimoRouter, nil
}

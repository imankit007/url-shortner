//go:build wireinject

package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/imankit007/url-shortner/controller"
	"github.com/imankit007/url-shortner/infrastructure"
	"github.com/imankit007/url-shortner/middleware"
	"github.com/imankit007/url-shortner/repository"
	"github.com/imankit007/url-shortner/service"
	"github.com/imankit007/url-shortner/utils/counter"
)

func InitializeRouter() (*gin.Engine, error) {
	panic(wire.Build(
		infrastructure.NewAuthServicePublicKey,
		infrastructure.NewAuthTokenIssuer,
		infrastructure.NewMongoClient,
		infrastructure.NewURLCollection,
		infrastructure.NewRedisClient,
		repository.NewMongoURLRepository,
		repository.NewCachedURLRepository,
		NewHashIDEncoder,
		NewApplicationBaseURL,
		counter.NewURLCodeCounter,
		service.NewTokenAuthenticationService,
		middleware.NewJWTAuthenticationMiddleware,
		service.NewURLService,
		controller.NewURLController,
		NewRouter,
	))
}

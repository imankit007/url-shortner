//go:build wireinject

package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/imankit007/url-shortner/controller"
	"github.com/imankit007/url-shortner/infrastructure"
	"github.com/imankit007/url-shortner/repository"
	"github.com/imankit007/url-shortner/service"
	"github.com/imankit007/url-shortner/utils/counter"
)

func InitializeRouter() (*gin.Engine, error) {
	panic(wire.Build(
		infrastructure.NewMongoClient,
		infrastructure.NewURLCollection,
		repository.NewMongoURLRepository,
		NewHashIDEncoder,
		NewApplicationBaseURL,
		counter.NewURLCodeCounter,
		service.NewURLService,
		controller.NewURLController,
		NewRouter,
	))
}

package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/imankit007/url-shortner-redirector-service/controller"
	"github.com/imankit007/url-shortner-redirector-service/infrastructure"
	"github.com/imankit007/url-shortner-redirector-service/repository"
	"github.com/imankit007/url-shortner-redirector-service/service"
	"github.com/speps/go-hashids"
)

func InitializeRouter() (*gin.Engine, error) {
	mongoClient, err := infrastructure.NewMongoClient()
	if err != nil {
		return nil, err
	}

	urlCollection := infrastructure.NewURLCollection(mongoClient)
	redisClient := infrastructure.NewRedisClient()
	hashIDEncoder, err := NewHashIDEncoder()
	if err != nil {
		return nil, err
	}

	mongoURLRepository := repository.NewMongoURLRepository(urlCollection)
	urlRepository := repository.NewCachedURLRepository(mongoURLRepository, redisClient)
	redirectService := service.NewRedirectService(urlRepository, hashIDEncoder)
	redirectController := controller.NewRedirectController(redirectService)

	engine := gin.Default()
	engine.GET("/:code", redirectController.RedirectToOriginalURLHandler)
	return engine, nil
}

func NewHashIDEncoder() (*hashids.HashID, error) {
	hashIDConfig := hashids.NewData()
	hashIDConfig.Salt = "my-secret-salt"
	hashIDConfig.MinLength = 6

	return hashids.NewWithData(hashIDConfig)
}

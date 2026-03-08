package bootstrap

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/imankit007/url-shortner/controller"
	"github.com/imankit007/url-shortner/middleware"
	"github.com/speps/go-hashids"
)

func NewHashIDEncoder() (*hashids.HashID, error) {
	hashIDConfig := hashids.NewData()
	hashIDConfig.Salt = "my-secret-salt"
	hashIDConfig.MinLength = 6

	return hashids.NewWithData(hashIDConfig)
}

func NewApplicationBaseURL() string {
	return "http://localhost:8080"
}

func NewRouter(
	urlController *controller.URLController,
	jwtAuthenticationMiddleware *middleware.JWTAuthenticationMiddleware,
) *gin.Engine {
	engine := gin.Default()
	engine.Use(cors.Default())

	apiV1 := engine.Group("/api/v1")
	apiV1.Use(jwtAuthenticationMiddleware.RequireAuthenticatedUser())
	apiV1.GET("/urls", urlController.ListURLMappingsHandler)
	apiV1.POST("/urls/shorten", urlController.CreateShortURLsHandler)
	engine.GET("/:code", urlController.RedirectToOriginalURLHandler)
	return engine
}

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imankit007/url-shortner/middleware"
	"github.com/imankit007/url-shortner/model"
	"github.com/imankit007/url-shortner/service"
)

type URLController struct {
	urlService service.URLService
}

func NewURLController(urlService service.URLService) *URLController {
	return &URLController{
		urlService: urlService,
	}
}

func (c *URLController) CreateShortURLsHandler(ctx *gin.Context) {
	authenticatedUser, ok := middleware.GetAuthenticatedUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return
	}

	var shortenRequest model.ShortenRequest
	if err := ctx.ShouldBindJSON(&shortenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	response, err := c.urlService.CreateShortURLs(ctx.Request.Context(), authenticatedUser, shortenRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save short url"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *URLController) ListURLMappingsHandler(ctx *gin.Context) {
	authenticatedUser, ok := middleware.GetAuthenticatedUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return
	}

	urlEntries, err := c.urlService.ListURLMappings(ctx.Request.Context(), authenticatedUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, urlEntries)
}

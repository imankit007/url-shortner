package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	var shortenRequest model.ShortenRequest
	if err := ctx.ShouldBindJSON(&shortenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	response, err := c.urlService.CreateShortURLs(ctx.Request.Context(), shortenRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save short url"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *URLController) RedirectToOriginalURLHandler(ctx *gin.Context) {
	redirectURL, err := c.urlService.ResolveRedirectURL(ctx.Request.Context(), ctx.Param("code"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidShortCode):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid short code"})
		case errors.Is(err, service.ErrShortURLNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "short url not found"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load short url"})
		}
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (c *URLController) ListURLMappingsHandler(ctx *gin.Context) {
	urlEntries, err := c.urlService.ListURLMappings(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, urlEntries)
}

package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imankit007/url-shortner-redirector-service/service"
)

type RedirectController struct {
	redirectService service.RedirectService
}

func NewRedirectController(redirectService service.RedirectService) *RedirectController {
	return &RedirectController{
		redirectService: redirectService,
	}
}

func (c *RedirectController) RedirectToOriginalURLHandler(ctx *gin.Context) {
	redirectURL, err := c.redirectService.ResolveRedirectURL(ctx.Request.Context(), ctx.Param("code"))
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

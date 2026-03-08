package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imankit007/url-shortner/model"
	"github.com/imankit007/url-shortner/service"
)

const AuthenticatedUserContextKey = "authenticated_user"

type JWTAuthenticationMiddleware struct {
	tokenAuthenticationService service.TokenAuthenticationService
}

func NewJWTAuthenticationMiddleware(tokenAuthenticationService service.TokenAuthenticationService) *JWTAuthenticationMiddleware {
	return &JWTAuthenticationMiddleware{
		tokenAuthenticationService: tokenAuthenticationService,
	}
}

func (m *JWTAuthenticationMiddleware) RequireAuthenticatedUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authorizationHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			ctx.Abort()
			return
		}

		authenticatedUser, err := m.tokenAuthenticationService.Authenticate(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			ctx.Abort()
			return
		}

		ctx.Set(AuthenticatedUserContextKey, authenticatedUser)
		ctx.Next()
	}
}

func GetAuthenticatedUser(ctx *gin.Context) (model.AuthenticatedUser, bool) {
	value, exists := ctx.Get(AuthenticatedUserContextKey)
	if !exists {
		return model.AuthenticatedUser{}, false
	}

	authenticatedUser, ok := value.(model.AuthenticatedUser)
	return authenticatedUser, ok
}

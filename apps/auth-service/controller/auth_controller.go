package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/imankit007/url-shortner-auth-service/model"
	"github.com/imankit007/url-shortner-auth-service/service"
)

type AuthController struct {
	userAuthenticationService service.UserAuthenticationService
	tokenService              service.TokenService
	publicKeyPEM              []byte
}

func NewAuthController(
	userAuthenticationService service.UserAuthenticationService,
	tokenService service.TokenService,
	publicKeyPEM []byte,
) *AuthController {
	return &AuthController{
		userAuthenticationService: userAuthenticationService,
		tokenService:              tokenService,
		publicKeyPEM:              publicKeyPEM,
	}
}

func (c *AuthController) CreateAccessTokenHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var tokenRequest model.TokenRequest
	if err := json.NewDecoder(request.Body).Decode(&tokenRequest); err != nil {
		writeJSON(responseWriter, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	authenticatedUser, err := c.userAuthenticationService.Authenticate(tokenRequest.Email, tokenRequest.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			writeJSON(responseWriter, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		default:
			writeJSON(responseWriter, http.StatusInternalServerError, map[string]string{"error": "failed to authenticate user"})
		}
		return
	}

	tokenResponse, err := c.tokenService.CreateToken(authenticatedUser)
	if err != nil {
		writeJSON(responseWriter, http.StatusInternalServerError, map[string]string{"error": "failed to create access token"})
		return
	}

	writeJSON(responseWriter, http.StatusOK, tokenResponse)
}

func (c *AuthController) PublicKeyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/x-pem-file")
	responseWriter.WriteHeader(http.StatusOK)
	_, _ = responseWriter.Write(c.publicKeyPEM)
}

func writeJSON(responseWriter http.ResponseWriter, statusCode int, payload any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	_ = json.NewEncoder(responseWriter).Encode(payload)
}

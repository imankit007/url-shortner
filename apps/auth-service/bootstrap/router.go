package bootstrap

import (
	"net/http"

	"github.com/imankit007/url-shortner-auth-service/controller"
	"github.com/imankit007/url-shortner-auth-service/infrastructure"
	"github.com/imankit007/url-shortner-auth-service/repository"
	"github.com/imankit007/url-shortner-auth-service/service"
)

func InitializeRouter() (http.Handler, error) {
	authPrivateKey, err := infrastructure.NewAuthServicePrivateKey()
	if err != nil {
		return nil, err
	}

	authPublicKeyPEM, err := infrastructure.NewAuthServicePublicKeyPEM()
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewInMemoryUserRepository()
	userAuthenticationService := service.NewUserAuthenticationService(userRepository)
	tokenService := service.NewTokenService(authPrivateKey, service.NewTokenIssuer(), service.NewTokenTTL())
	authController := controller.NewAuthController(userAuthenticationService, tokenService, authPublicKeyPEM)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/token", authController.CreateAccessTokenHandler)
	mux.HandleFunc("GET /api/v1/auth/public-key", authController.PublicKeyHandler)
	return mux, nil
}

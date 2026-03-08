package service

import (
	"crypto/rsa"
	"time"

	"github.com/imankit007/url-shortner-auth-service/model"
	"github.com/imankit007/url-shortner-auth-service/security"
)

type TokenService interface {
	CreateToken(authenticatedUser model.AuthenticatedUser) (model.TokenResponse, error)
}

type tokenService struct {
	privateKey *rsa.PrivateKey
	issuer     string
	tokenTTL   time.Duration
}

func NewTokenService(privateKey *rsa.PrivateKey, issuer string, tokenTTL time.Duration) TokenService {
	return &tokenService{
		privateKey: privateKey,
		issuer:     issuer,
		tokenTTL:   tokenTTL,
	}
}

func NewTokenIssuer() string {
	return "url-shortener-auth-service"
}

func NewTokenTTL() time.Duration {
	return time.Hour
}

func (s *tokenService) CreateToken(authenticatedUser model.AuthenticatedUser) (model.TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(s.tokenTTL)

	token, err := security.GenerateRS256Token(s.privateKey, security.TokenClaims{
		Issuer:    s.issuer,
		Subject:   authenticatedUser.UserID,
		UserID:    authenticatedUser.UserID,
		TenantID:  authenticatedUser.TenantID,
		Email:     authenticatedUser.Email,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	})
	if err != nil {
		return model.TokenResponse{}, err
	}

	return model.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.Unix(),
		UserID:      authenticatedUser.UserID,
		TenantID:    authenticatedUser.TenantID,
		Email:       authenticatedUser.Email,
	}, nil
}

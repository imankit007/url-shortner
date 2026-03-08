package service

import (
	"crypto/rsa"
	"errors"

	"github.com/imankit007/url-shortner/model"
	"github.com/imankit007/url-shortner/security"
)

var ErrInvalidAccessToken = errors.New("invalid access token")

type TokenAuthenticationService interface {
	Authenticate(accessToken string) (model.AuthenticatedUser, error)
}

type tokenAuthenticationService struct {
	authServicePublicKey *rsa.PublicKey
	expectedIssuer       string
}

func NewTokenAuthenticationService(
	authServicePublicKey *rsa.PublicKey,
	expectedIssuer string,
) TokenAuthenticationService {
	return &tokenAuthenticationService{
		authServicePublicKey: authServicePublicKey,
		expectedIssuer:       expectedIssuer,
	}
}

func (s *tokenAuthenticationService) Authenticate(accessToken string) (model.AuthenticatedUser, error) {
	claims, err := security.VerifyRS256Token(accessToken, s.authServicePublicKey, s.expectedIssuer)
	if err != nil {
		return model.AuthenticatedUser{}, ErrInvalidAccessToken
	}

	return model.AuthenticatedUser{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Email:    claims.Email,
	}, nil
}

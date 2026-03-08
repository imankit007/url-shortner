package service

import (
	"errors"

	"github.com/imankit007/url-shortner-auth-service/model"
	"github.com/imankit007/url-shortner-auth-service/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserAuthenticationService interface {
	Authenticate(email string, password string) (model.AuthenticatedUser, error)
}

type userAuthenticationService struct {
	userRepository repository.UserRepository
}

func NewUserAuthenticationService(userRepository repository.UserRepository) UserAuthenticationService {
	return &userAuthenticationService{
		userRepository: userRepository,
	}
}

func (s *userAuthenticationService) Authenticate(email string, password string) (model.AuthenticatedUser, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return model.AuthenticatedUser{}, ErrInvalidCredentials
	}

	if user.PasswordHash != password {
		return model.AuthenticatedUser{}, ErrInvalidCredentials
	}

	return model.AuthenticatedUser{
		UserID:   user.UserID,
		TenantID: user.TenantID,
		Email:    user.Email,
	}, nil
}

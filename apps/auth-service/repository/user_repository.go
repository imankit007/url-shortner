package repository

import (
	"errors"

	"github.com/imankit007/url-shortner-auth-service/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	FindByEmail(email string) (model.User, error)
}

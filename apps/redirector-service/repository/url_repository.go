package repository

import (
	"context"
	"errors"

	"github.com/imankit007/url-shortner-redirector-service/model"
)

var ErrURLNotFound = errors.New("url not found")

type URLRepository interface {
	FindByCode(ctx context.Context, code int64) (model.URLEntry, error)
}

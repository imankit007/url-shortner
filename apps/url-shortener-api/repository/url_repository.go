package repository

import (
	"context"
	"errors"

	"github.com/imankit007/url-shortner/model"
)

var ErrURLNotFound = errors.New("url not found")

type URLRepository interface {
	Save(ctx context.Context, entry model.URLEntry) error
	FindByCode(ctx context.Context, code int64) (model.URLEntry, error)
	ListAll(ctx context.Context) ([]model.URLEntry, error)
}

package service

import (
	"context"
	"errors"
	"strings"

	"github.com/imankit007/url-shortner-redirector-service/repository"
	"github.com/speps/go-hashids"
)

var (
	ErrInvalidShortCode = errors.New("invalid short code")
	ErrShortURLNotFound = errors.New("short url not found")
)

type RedirectService interface {
	ResolveRedirectURL(ctx context.Context, code string) (string, error)
}

type redirectService struct {
	urlRepository repository.URLRepository
	hashIDEncoder *hashids.HashID
}

func NewRedirectService(urlRepository repository.URLRepository, hashIDEncoder *hashids.HashID) RedirectService {
	return &redirectService{
		urlRepository: urlRepository,
		hashIDEncoder: hashIDEncoder,
	}
}

func (s *redirectService) ResolveRedirectURL(ctx context.Context, code string) (string, error) {
	decodedValues, err := s.hashIDEncoder.DecodeInt64WithError(code)
	if err != nil || len(decodedValues) == 0 {
		return "", ErrInvalidShortCode
	}

	urlEntry, err := s.urlRepository.FindByCode(ctx, decodedValues[0])
	if err != nil {
		if errors.Is(err, repository.ErrURLNotFound) {
			return "", ErrShortURLNotFound
		}
		return "", err
	}

	if strings.HasPrefix(urlEntry.OriginalURL, "http://") || strings.HasPrefix(urlEntry.OriginalURL, "https://") {
		return urlEntry.OriginalURL, nil
	}

	return "https://" + urlEntry.OriginalURL, nil
}

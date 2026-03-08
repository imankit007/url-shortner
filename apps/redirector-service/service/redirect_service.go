package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/imankit007/url-shortner-redirector-service/model"
	"github.com/imankit007/url-shortner-redirector-service/repository"
	"github.com/speps/go-hashids"
)

var (
	ErrInvalidShortCode = errors.New("invalid short code")
	ErrShortURLNotFound = errors.New("short url not found")
)

type RedirectRequest struct {
	Code      string
	UserAgent string
	Referer   string
	IPAddress string
}

type RedirectService interface {
	ResolveRedirectURL(ctx context.Context, req RedirectRequest) (string, error)
}

type redirectService struct {
	urlRepository       repository.URLRepository
	hashIDEncoder       *hashids.HashID
	clickEventPublisher ClickEventPublisher
}

func NewRedirectService(
	urlRepository repository.URLRepository,
	hashIDEncoder *hashids.HashID,
	clickEventPublisher ClickEventPublisher,
) RedirectService {
	return &redirectService{
		urlRepository:       urlRepository,
		hashIDEncoder:       hashIDEncoder,
		clickEventPublisher: clickEventPublisher,
	}
}

func (s *redirectService) ResolveRedirectURL(ctx context.Context, req RedirectRequest) (string, error) {
	decodedValues, err := s.hashIDEncoder.DecodeInt64WithError(req.Code)

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

	redirectURL := urlEntry.OriginalURL
	if !strings.HasPrefix(redirectURL, "http://") && !strings.HasPrefix(redirectURL, "https://") {
		redirectURL = "https://" + redirectURL
	}

	s.clickEventPublisher.Publish(ctx, model.ClickEvent{
		ShortCode:   req.Code,
		OriginalURL: urlEntry.OriginalURL,
		TenantID:    urlEntry.TenantID,
		Timestamp:   time.Now().UTC(),
		UserAgent:   req.UserAgent,
		Referer:     req.Referer,
		IPAddress:   req.IPAddress,
	})

	return redirectURL, nil
}

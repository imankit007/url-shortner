package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/imankit007/url-shortner/model"
	"github.com/imankit007/url-shortner/repository"
	"github.com/imankit007/url-shortner/utils/counter"
	"github.com/speps/go-hashids"
)

var (
	ErrInvalidShortCode = errors.New("invalid short code")
	ErrShortURLNotFound = errors.New("short url not found")
	ErrInvalidURL       = errors.New("invalid url")
)

type URLService interface {
	CreateShortURLs(ctx context.Context, authenticatedUser model.AuthenticatedUser, shortenRequest model.ShortenRequest) ([]model.ShortenResponse, error)
	ListURLMappings(ctx context.Context, authenticatedUser model.AuthenticatedUser) ([]model.URLEntry, error)
}

type urlService struct {
	urlRepository  repository.URLRepository
	hashIDEncoder  *hashids.HashID
	urlCodeCounter *counter.URLCodeCounter
	baseURL        string
}

func NewURLService(
	urlRepository repository.URLRepository,
	hashIDEncoder *hashids.HashID,
	urlCodeCounter *counter.URLCodeCounter,
	baseURL string,
) URLService {
	return &urlService{
		urlRepository:  urlRepository,
		hashIDEncoder:  hashIDEncoder,
		urlCodeCounter: urlCodeCounter,
		baseURL:        strings.TrimRight(baseURL, "/"),
	}
}

func (s *urlService) CreateShortURLs(
	ctx context.Context,
	authenticatedUser model.AuthenticatedUser,
	shortenRequest model.ShortenRequest,
) ([]model.ShortenResponse, error) {
	shortenResponses := make([]model.ShortenResponse, 0, len(shortenRequest.Links))
	for _, linkRequest := range shortenRequest.Links {
		if err := validateURL(linkRequest.URL); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidURL, linkRequest.URL)
		}

		nextCodeValue, err := s.urlCodeCounter.Next(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		shortCode, err := s.hashIDEncoder.EncodeInt64([]int64{int64(nextCodeValue)})
		if err != nil {
			return nil, err
		}

		urlEntry := model.URLEntry{
			OriginalURL:     linkRequest.URL,
			Code:            int64(nextCodeValue),
			TenantID:        authenticatedUser.TenantID,
			CreatedByUserID: authenticatedUser.UserID,
		}

		if err := s.urlRepository.Save(ctx, urlEntry); err != nil {
			return nil, err
		}

		shortenResponses = append(shortenResponses, model.ShortenResponse{
			LongURL:  linkRequest.URL,
			ShortURL: s.baseURL + "/" + shortCode,
		})
	}

	return shortenResponses, nil
}

func (s *urlService) ListURLMappings(ctx context.Context, authenticatedUser model.AuthenticatedUser) ([]model.URLEntry, error) {
	return s.urlRepository.ListAllByTenant(ctx, authenticatedUser.TenantID)
}

func validateURL(rawURL string) error {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("unsupported scheme: %s", parsed.Scheme)
	}

	if parsed.Host == "" {
		return fmt.Errorf("missing host")
	}

	return nil
}

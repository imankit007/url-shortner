package service

import (
	"context"
	"errors"
	"strings"

	"github.com/imankit007/url-shortner/model"
	"github.com/imankit007/url-shortner/repository"
	"github.com/imankit007/url-shortner/utils/counter"
	"github.com/speps/go-hashids"
)

var (
	ErrInvalidShortCode = errors.New("invalid short code")
	ErrShortURLNotFound = errors.New("short url not found")
)

type URLService interface {
	CreateShortURLs(ctx context.Context, authenticatedUser model.AuthenticatedUser, shortenRequest model.ShortenRequest) ([]model.ShortenResponse, error)
	ResolveRedirectURL(ctx context.Context, code string) (string, error)
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
		nextCodeValue := s.urlCodeCounter.Next()

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

func (s *urlService) ResolveRedirectURL(ctx context.Context, code string) (string, error) {
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

func (s *urlService) ListURLMappings(ctx context.Context, authenticatedUser model.AuthenticatedUser) ([]model.URLEntry, error) {
	return s.urlRepository.ListAllByTenant(ctx, authenticatedUser.TenantID)
}

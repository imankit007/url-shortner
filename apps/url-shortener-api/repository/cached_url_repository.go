package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/imankit007/url-shortner/model"
	"github.com/redis/go-redis/v9"
)

const shortCodeCacheTTL = 24 * time.Hour

type CachedURLRepository struct {
	baseRepository URLStorageRepository
	redisClient    *redis.Client
}

func NewCachedURLRepository(baseRepository URLStorageRepository, redisClient *redis.Client) URLRepository {
	return &CachedURLRepository{
		baseRepository: baseRepository,
		redisClient:    redisClient,
	}
}

func (r *CachedURLRepository) Save(ctx context.Context, entry model.URLEntry) error {
	if err := r.baseRepository.Save(ctx, entry); err != nil {
		return err
	}

	r.cacheEntry(ctx, entry)
	return nil
}

func (r *CachedURLRepository) FindByCode(ctx context.Context, code int64) (model.URLEntry, error) {
	if entry, found := r.loadFromCache(ctx, code); found {
		return entry, nil
	}

	entry, err := r.baseRepository.FindByCode(ctx, code)
	if err != nil {
		return model.URLEntry{}, err
	}

	r.cacheEntry(ctx, entry)
	return entry, nil
}

func (r *CachedURLRepository) ListAllByTenant(ctx context.Context, tenantID string) ([]model.URLEntry, error) {
	return r.baseRepository.ListAllByTenant(ctx, tenantID)
}

func (r *CachedURLRepository) loadFromCache(ctx context.Context, code int64) (model.URLEntry, bool) {
	if r.redisClient == nil {
		return model.URLEntry{}, false
	}

	cachedValue, err := r.redisClient.Get(ctx, r.cacheKey(code)).Result()
	if err != nil {
		return model.URLEntry{}, false
	}

	var entry model.URLEntry
	if err := json.Unmarshal([]byte(cachedValue), &entry); err != nil {
		return model.URLEntry{}, false
	}

	return entry, true
}

func (r *CachedURLRepository) cacheEntry(ctx context.Context, entry model.URLEntry) {
	if r.redisClient == nil {
		return
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return
	}

	_ = r.redisClient.Set(ctx, r.cacheKey(entry.Code), payload, shortCodeCacheTTL).Err()
}

func (r *CachedURLRepository) cacheKey(code int64) string {
	return fmt.Sprintf("short-url:code:%d", code)
}

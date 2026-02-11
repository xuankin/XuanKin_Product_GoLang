package repository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
}
type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) CacheRepository {
	return &cacheRepository{client: client}
}
func (repo *cacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := repo.client.Get(ctx, key).Result()
	if err != nil {
		return err

	}
	return json.Unmarshal([]byte(val), dest)
}
func (repo *cacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = repo.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
func (repo *cacheRepository) Delete(ctx context.Context, key string) error {
	return repo.client.Del(ctx, key).Err()
}
func (repo *cacheRepository) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := repo.client.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			repo.client.Del(ctx, keys...)
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

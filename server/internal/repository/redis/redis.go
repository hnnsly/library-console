package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hnnsly/library-console/internal/config" // Обновите путь
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
	TTL    time.Duration
}

func New(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	ttl := time.Duration(cfg.CacheTTLSeconds) * time.Second
	if cfg.CacheTTLSeconds == 0 {
		ttl = 5 * time.Minute // Default TTL if not configured
	}

	return &Redis{Client: rdb, TTL: ttl}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil // Key does not exist
	} else if err != nil {
		return fmt.Errorf("redis get error for key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal redis value for key %s: %w", key, err)
	}
	return nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, ttlOverride ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for redis key %s: %w", key, err)
	}

	ttl := r.TTL
	if len(ttlOverride) > 0 {
		ttl = ttlOverride[0]
	}

	if err := r.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set error for key %s: %w", key, err)
	}
	return nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete error for key %s: %w", key, err)
	}
	return nil
}

// Пример специфичных для домена методов кеширования, как в aqua-taxi
// Вы можете добавить больше таких методов в pharmacy/repository/repository.go

func (r *Redis) CachePopularProducts(ctx context.Context, products interface{}) error {
	return r.Set(ctx, "popular_products", products)
}

func (r *Redis) GetPopularProducts(ctx context.Context) (interface{}, error) {
	var products []interface{} // Или конкретный тип
	err := r.Get(ctx, "popular_products", &products)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 { // redis.Nil обработан в r.Get, это для случая пустого массива
		return nil, nil
	}
	return products, nil
}

package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hnnsly/library-console/internal/config"
	"github.com/redis/go-redis/v9"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrKeyExists   = errors.New("key already exists")
)

type Redis struct {
	Client *redis.Client
	TTL    time.Duration
}

func New(ctx context.Context, cfg config.RedisConfig) (*Redis, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("redis address is required")
	}

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

// Get получает значение по ключу и десериализует его в dest
func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if dest == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrKeyNotFound
	} else if err != nil {
		return fmt.Errorf("redis get error for key %s: %w", key, err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal redis value for key %s: %w", key, err)
	}
	return nil
}

// Set устанавливает значение по ключу с TTL
func (r *Redis) Set(ctx context.Context, key string, value interface{}, ttlOverride ...time.Duration) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if value == nil {
		return fmt.Errorf("value cannot be nil")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for redis key %s: %w", key, err)
	}

	ttl := r.TTL
	if len(ttlOverride) > 0 && ttlOverride[0] > 0 {
		ttl = ttlOverride[0]
	}

	if err := r.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set error for key %s: %w", key, err)
	}
	return nil
}

// SetNX устанавливает значение только если ключ не существует
func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, ttlOverride ...time.Duration) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if value == nil {
		return fmt.Errorf("value cannot be nil")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for redis key %s: %w", key, err)
	}

	ttl := r.TTL
	if len(ttlOverride) > 0 && ttlOverride[0] > 0 {
		ttl = ttlOverride[0]
	}

	result, err := r.Client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		return fmt.Errorf("redis setnx error for key %s: %w", key, err)
	}
	if !result {
		return ErrKeyExists
	}
	return nil
}

// Del удаляет ключ
func (r *Redis) Del(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if err := r.Client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis delete error for key %s: %w", key, err)
	}
	return nil
}

// Exists проверяет существование ключа
func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key cannot be empty")
	}

	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists error for key %s: %w", key, err)
	}
	return result > 0, nil
}

// Expire устанавливает TTL для существующего ключа
func (r *Redis) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if ttl <= 0 {
		return fmt.Errorf("ttl must be positive")
	}

	if err := r.Client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("redis expire error for key %s: %w", key, err)
	}
	return nil
}

// TTLRemaining возвращает оставшееся время жизни ключа
func (r *Redis) TTLRemaining(ctx context.Context, key string) (time.Duration, error) {
	if key == "" {
		return 0, fmt.Errorf("key cannot be empty")
	}

	ttl, err := r.Client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis ttl error for key %s: %w", key, err)
	}
	return ttl, nil
}

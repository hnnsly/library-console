package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hnnsly/library-console/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Redis struct {
	conn *redis.Client
}

func New(ctx context.Context, config config.Redis) (*Redis, error) {
	log.Info().Msgf("Initializing Redis connections to %s:%s", config.Host, config.Port)

	сonn := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           0,
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		PoolSize:     10,
		MinIdleConns: 5,
		PoolTimeout:  300 * time.Millisecond,
	})
	db := &Redis{
		conn: сonn,
	}

	if err := db.TestConn(ctx); err != nil {
		return nil, err
	}

	log.Info().Msg("Redis connections initialized successfully")
	return db, nil
}

func (db *Redis) TestConn(ctx context.Context) error {
	if _, err := db.conn.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to main Redis: %w", err)
	}

	return nil
}

func (rd *Redis) Close() error {
	if err := rd.conn.Close(); err != nil {
		return err
	}

	log.Warn().Msg("Redis repository is successfully closed")
	return nil
}

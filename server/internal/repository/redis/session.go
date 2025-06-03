package redis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	sessionKeyPrefix = "session:"
)

// SessionValue - минимальная информация, хранимая в сессии
type SessionValue struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

// CreateSession создает новую сессию для пользователя
func (r *Redis) CreateSession(ctx context.Context, userID int64, role string, ttl time.Duration) (string, error) {
	// Генерируем случайный session ID
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("ошибка генерации session ID: %w", err)
	}
	sessionID := base64.RawURLEncoding.EncodeToString(randomBytes)

	key := sessionKeyPrefix + sessionID

	sessionData := SessionValue{
		UserID: userID,
		Role:   role,
	}

	value, err := json.Marshal(sessionData)
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации данных сессии: %w", err)
	}

	err = r.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return "", fmt.Errorf("не удалось создать сессию: %w", err)
	}

	return sessionID, nil
}

// GetSession получает данные сессии по ID
func (r *Redis) GetSession(ctx context.Context, sessionID string) (*SessionValue, error) {
	key := sessionKeyPrefix + sessionID

	value, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("сессия не найдена")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении сессии: %w", err)
	}

	var sessionData SessionValue
	if err := json.Unmarshal([]byte(value), &sessionData); err != nil {
		return nil, fmt.Errorf("ошибка десериализации данных сессии: %w", err)
	}

	return &sessionData, nil
}

// RefreshSession обновляет TTL сессии
func (r *Redis) RefreshSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	key := sessionKeyPrefix + sessionID

	exists, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования сессии: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("сессия не найдена")
	}

	success, err := r.Client.Expire(ctx, key, ttl).Result()
	if err != nil {
		return fmt.Errorf("ошибка при обновлении TTL сессии: %w", err)
	}
	if !success {
		return fmt.Errorf("не удалось обновить TTL сессии")
	}

	return nil
}

// DeleteSession удаляет сессию
func (r *Redis) DeleteSession(ctx context.Context, sessionID string) error {
	key := sessionKeyPrefix + sessionID
	_, err := r.Client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("ошибка при удалении сессии: %w", err)
	}
	return nil
}

// DeleteAllUserSessions удаляет все сессии пользователя
func (r *Redis) DeleteAllUserSessions(ctx context.Context, userID int64) error {
	pattern := sessionKeyPrefix + "*"
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()

	var keysToDelete []string
	for iter.Next(ctx) {
		key := iter.Val()

		value, err := r.Client.Get(ctx, key).Result()
		if err != nil {
			continue // Пропускаем ошибочные ключи
		}

		var sessionData SessionValue
		if err := json.Unmarshal([]byte(value), &sessionData); err != nil {
			continue
		}

		if sessionData.UserID == userID {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("ошибка при сканировании сессий: %w", err)
	}

	if len(keysToDelete) > 0 {
		_, err := r.Client.Del(ctx, keysToDelete...).Result()
		if err != nil {
			return fmt.Errorf("ошибка при удалении сессий пользователя: %w", err)
		}
	}

	return nil
}

// DeleteAllUserSessionsExcept удаляет все сессии пользователя кроме указанной
func (r *Redis) DeleteAllUserSessionsExcept(ctx context.Context, userID int64, exceptSessionID string) error {
	pattern := sessionKeyPrefix + "*"
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()

	var keysToDelete []string
	for iter.Next(ctx) {
		key := iter.Val()
		currentSessionID := key[len(sessionKeyPrefix):] // Убираем префикс

		if currentSessionID == exceptSessionID {
			continue
		}

		value, err := r.Client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var sessionData SessionValue
		if err := json.Unmarshal([]byte(value), &sessionData); err != nil {
			continue
		}

		if sessionData.UserID == userID {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("ошибка при сканировании сессий: %w", err)
	}

	if len(keysToDelete) > 0 {
		_, err := r.Client.Del(ctx, keysToDelete...).Result()
		if err != nil {
			return fmt.Errorf("ошибка при удалении сессий пользователя: %w", err)
		}
	}

	return nil
}

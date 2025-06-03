package redis

import (
	"context"
	"fmt"
	"time"
)

// CreateSession создает сессию с проверкой на существование
func (r *Redis) CreateSession(ctx context.Context, sessionID string, data interface{}, ttl time.Duration) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	return r.SetNX(ctx, sessionKey, data, ttl)
}

// GetSession получает данные сессии
func (r *Redis) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	return r.Get(ctx, sessionKey, dest)
}

// UpdateSessionTTL продлевает сессию
func (r *Redis) UpdateSessionTTL(ctx context.Context, sessionID string, ttl time.Duration) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)

	// Проверяем, что сессия существует
	exists, err := r.Exists(ctx, sessionKey)
	if err != nil {
		return err
	}
	if !exists {
		return ErrKeyNotFound
	}

	return r.Expire(ctx, sessionKey, ttl)
}

// DeleteSession удаляет сессию
func (r *Redis) DeleteSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	return r.Del(ctx, sessionKey)
}

// SessionExists проверяет существование сессии
func (r *Redis) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	if sessionID == "" {
		return false, fmt.Errorf("session ID cannot be empty")
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	return r.Exists(ctx, sessionKey)
}

// GetSessionTTL возвращает оставшееся время жизни сессии
func (r *Redis) GetSessionTTL(ctx context.Context, sessionID string) (time.Duration, error) {
	if sessionID == "" {
		return 0, fmt.Errorf("session ID cannot be empty")
	}

	sessionKey := fmt.Sprintf("session:%s", sessionID)
	return r.TTLRemaining(ctx, sessionKey)
}

// InvalidateAllUserSessions инвалидирует все сессии пользователя по паттерну
func (r *Redis) InvalidateAllUserSessions(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return fmt.Errorf("user ID must be positive")
	}

	// Используем SCAN для поиска всех ключей сессий
	pattern := "session:*"
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()

	var sessionsToDelete []string

	for iter.Next(ctx) {
		key := iter.Val()

		// Получаем данные сессии для проверки userID
		var sessionData struct {
			UserID int64 `json:"user_id"`
		}

		if err := r.Get(ctx, key, &sessionData); err != nil {
			// Пропускаем некорректные сессии
			continue
		}

		if sessionData.UserID == userID {
			sessionsToDelete = append(sessionsToDelete, key)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan sessions: %w", err)
	}

	// Удаляем найденные сессии
	for _, key := range sessionsToDelete {
		if err := r.Del(ctx, key); err != nil {
			// Логируем ошибку, но продолжаем удаление других сессий
			fmt.Printf("Warning: failed to delete session %s: %v\n", key, err)
		}
	}

	return nil
}

// GetActiveSessionsCount возвращает количество активных сессий
func (r *Redis) GetActiveSessionsCount(ctx context.Context) (int64, error) {
	pattern := "session:*"
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()

	var count int64
	for iter.Next(ctx) {
		count++
	}

	if err := iter.Err(); err != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", err)
	}

	return count, nil
}

// CleanupExpiredSessions удаляет истекшие сессии (Redis делает это автоматически, но может быть полезно для статистики)
func (r *Redis) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	pattern := "session:*"
	iter := r.Client.Scan(ctx, 0, pattern, 0).Iterator()

	var deletedCount int64

	for iter.Next(ctx) {
		key := iter.Val()

		// Проверяем TTL
		ttl, err := r.TTLRemaining(ctx, key)
		if err != nil {
			continue
		}

		// Если TTL < 0, ключ уже истек (но Redis еще не удалил)
		if ttl < 0 {
			if err := r.Del(ctx, key); err == nil {
				deletedCount++
			}
		}
	}

	if err := iter.Err(); err != nil {
		return deletedCount, fmt.Errorf("failed to cleanup sessions: %w", err)
	}

	return deletedCount, nil
}

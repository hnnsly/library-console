package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/redis"
)

type SessionData struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

type SessionManager struct {
	redis *redis.Redis
}

func NewSessionManager(redis *redis.Redis) *SessionManager {
	return &SessionManager{
		redis: redis,
	}
}

func (sm *SessionManager) CreateSession(ctx context.Context, data SessionData) (string, error) {
	sessionID := uuid.New().String()
	sessionKey := fmt.Sprintf("session:%s", sessionID)

	sessionJSON, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session data: %w", err)
	}

	// Сессия действует 24 часа
	err = sm.redis.Set(ctx, sessionKey, sessionJSON, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to save session: %w", err)
	}

	return sessionID, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	sessionKey := fmt.Sprintf("session:%s", sessionID)

	var data SessionData
	err := sm.redis.Get(ctx, sessionKey, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Проверяем, что сессия не слишком старая (дополнительная защита)
	if time.Since(data.CreatedAt) > 7*24*time.Hour {
		sm.DeleteSession(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return &data, nil
}

func (sm *SessionManager) UpdateLastSeen(ctx context.Context, sessionID string) error {
	sessionKey := fmt.Sprintf("session:%s", sessionID)

	var data SessionData
	err := sm.redis.Get(ctx, sessionKey, &data)
	if err != nil {
		return fmt.Errorf("failed to get session for update: %w", err)
	}

	data.LastSeen = time.Now()

	sessionJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal updated session data: %w", err)
	}

	// Продлеваем сессию на 24 часа
	err = sm.redis.Set(ctx, sessionKey, sessionJSON, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	return sm.redis.Del(ctx, sessionKey)
}

func (sm *SessionManager) InvalidateUserSessions(ctx context.Context, userID int64) error {
	return sm.redis.InvalidateAllUserSessions(ctx, userID)
}

func (sm *SessionManager) InvalidateUserSessionsExcept(ctx context.Context, userID int64, exceptSessionID string) error {
	// Получаем все сессии пользователя и удаляем все кроме указанной
	pattern := "session:*"
	iter := sm.redis.Client.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		sessionID := key[8:] // убираем префикс "session:"

		if sessionID == exceptSessionID {
			continue
		}

		var sessionData SessionData
		if err := sm.redis.Get(ctx, key, &sessionData); err != nil {
			continue
		}

		if sessionData.UserID == userID {
			sm.redis.Del(ctx, key)
		}
	}

	return iter.Err()
}

func (sm *SessionManager) ValidateSessionSecurity(ctx context.Context, sessionID, currentIP, currentUserAgent string) error {
	session, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Проверяем IP (опционально, может быть проблематично с мобильными сетями)
	if session.IPAddress != currentIP {
		// Можно добавить логирование подозрительной активности
		// или требовать повторную аутентификацию
	}

	// Проверяем User-Agent (базовая проверка)
	if session.UserAgent != currentUserAgent {
		// Логируем потенциальную проблему безопасности
	}

	return nil
}

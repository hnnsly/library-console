package auth

import (
	"context"
	"time"

	"github.com/hnnsly/library-console/internal/repository/redis"
)

type SessionManager struct {
	redis *redis.Redis
}

func NewSessionManager(redis *redis.Redis) *SessionManager {
	return &SessionManager{
		redis: redis,
	}
}

// CreateSession создает новую сессию
func (sm *SessionManager) CreateSession(ctx context.Context, userID int64, role string, ttl time.Duration) (string, error) {
	return sm.redis.CreateSession(ctx, userID, role, ttl)
}

// GetSession получает данные сессии
func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*redis.SessionValue, error) {
	return sm.redis.GetSession(ctx, sessionID)
}

// RefreshSession обновляет TTL сессии
func (sm *SessionManager) RefreshSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	return sm.redis.RefreshSession(ctx, sessionID, ttl)
}

// DeleteSession удаляет сессию
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	return sm.redis.DeleteSession(ctx, sessionID)
}

// InvalidateAllUserSessions удаляет все сессии пользователя
func (sm *SessionManager) InvalidateAllUserSessions(ctx context.Context, userID int64) error {
	return sm.redis.DeleteAllUserSessions(ctx, userID)
}

// InvalidateAllUserSessionsExcept удаляет все сессии пользователя кроме указанной
func (sm *SessionManager) InvalidateAllUserSessionsExcept(ctx context.Context, userID int64, exceptSessionID string) error {
	return sm.redis.DeleteAllUserSessionsExcept(ctx, userID, exceptSessionID)
}

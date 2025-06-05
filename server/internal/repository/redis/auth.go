package redis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
)

const (
	sessionKeyPrefix       = "session:"
	emailCodeKeyPrefix     = "check:email:"
	verifiedEmailKeyPrefix = "verified:email:" // Новый префикс
)

type Session struct {
	ID     string
	UserID uuid.UUID
	Role   postgres.UserRole
}

// Создает новую сессию для пользователя
func (r *Redis) CreateSession(ctx context.Context, userID uuid.UUID, role postgres.UserRole, ttl time.Duration) (string, error) {
	// Генерируем ID сессии
	randomBytes := make([]byte, 20)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("ошибка генерации случайных данных: %w", err)
	}
	sessionID := base64.RawURLEncoding.EncodeToString(randomBytes)

	key := sessionKeyPrefix + sessionID

	// Создаем структуру сессии
	session := Session{
		ID:     sessionID,
		UserID: userID,
		Role:   role,
	}

	// Сохраняем структуру Session как hash
	sessionData := map[string]interface{}{
		"id":     session.ID,
		"userID": session.UserID.String(),
		"role":   string(session.Role),
	}

	// Используем HMSet для установки всех полей hash одновременно
	err := r.conn.HMSet(ctx, key, sessionData).Err()
	if err != nil {
		return "", fmt.Errorf("не удалось создать сессию: %w", err)
	}

	// Устанавливаем TTL для ключа
	err = r.conn.Expire(ctx, key, ttl).Err()
	if err != nil {
		return "", fmt.Errorf("не удалось установить TTL для сессии: %w", err)
	}

	return sessionID, nil
}

// GetSession получает сессию из Redis по ID
func (r *Redis) GetSession(ctx context.Context, sessionID string) (Session, error) {
	key := sessionKeyPrefix + sessionID

	// Получаем все поля hash
	sessionData, err := r.conn.HGetAll(ctx, key).Result()
	if err != nil {
		return Session{}, fmt.Errorf("ошибка при получении сессии: %w", err)
	}

	// Проверяем, что hash существует и не пуст
	if len(sessionData) == 0 {
		return Session{}, fmt.Errorf("сессия не найдена")
	}

	// Парсим UUID
	userID, err := uuid.Parse(sessionData["userID"])
	if err != nil {
		return Session{}, fmt.Errorf("ошибка парсинга UserID: %w", err)
	}

	session := Session{
		ID:     sessionData["id"],
		UserID: userID,
		Role:   postgres.UserRole(sessionData["role"]),
	}

	return session, nil
}

func (r *Redis) TerminateOtherSessions(ctx context.Context, userID uuid.UUID, currentSessionID string) error {
	pattern := sessionKeyPrefix + "*"

	// Получаем все ключи сессий
	keys, err := r.conn.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("ошибка при получении списка сессий: %w", err)
	}

	var keysToDelete []string

	// Проверяем каждую сессию
	for _, key := range keys {
		// Извлекаем sessionID из ключа
		sessionID := key[len(sessionKeyPrefix):]

		// Пропускаем текущую сессию
		if sessionID == currentSessionID {
			continue
		}

		// Получаем данные сессии
		sessionData, err := r.conn.HGetAll(ctx, key).Result()
		if err != nil {
			// Если не можем получить данные сессии, пропускаем
			continue
		}

		// Проверяем, что сессия не пуста
		if len(sessionData) == 0 {
			continue
		}

		// Проверяем, принадлежит ли сессия этому пользователю
		storedUserID, err := uuid.Parse(sessionData["userID"])
		if err != nil {
			// Если не можем распарсить UUID, пропускаем
			continue
		}

		if storedUserID == userID {
			keysToDelete = append(keysToDelete, key)
		}
	}

	// Удаляем найденные сессии
	if len(keysToDelete) > 0 {
		_, err := r.conn.Del(ctx, keysToDelete...).Result()
		if err != nil {
			return fmt.Errorf("ошибка при удалении сессий: %w", err)
		}
	}

	return nil
}

// GetUserSessions возвращает все активные сессии для указанного пользователя
func (r *Redis) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	pattern := sessionKeyPrefix + "*"

	// Получаем все ключи сессий
	keys, err := r.conn.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка сессий: %w", err)
	}

	var userSessions []Session

	// Проверяем каждую сессию
	for _, key := range keys {
		// Получаем данные сессии
		sessionData, err := r.conn.HGetAll(ctx, key).Result()
		if err != nil {
			// Если не можем получить данные сессии, пропускаем
			continue
		}

		// Проверяем, что сессия не пуста
		if len(sessionData) == 0 {
			continue
		}

		// Проверяем, принадлежит ли сессия этому пользователю
		storedUserID, err := uuid.Parse(sessionData["userID"])
		if err != nil {
			// Если не можем распарсить UUID, пропускаем
			continue
		}

		if storedUserID == userID {
			// Создаем объект сессии
			session := Session{
				ID:     sessionData["id"],
				UserID: storedUserID,
				Role:   postgres.UserRole(sessionData["role"]),
			}
			userSessions = append(userSessions, session)
		}
	}

	return userSessions, nil
}

// Удаляет сессию
func (r *Redis) DeleteSession(ctx context.Context, sessionID string) error {
	key := sessionKeyPrefix + sessionID
	_, err := r.conn.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("ошибка при удалении сессии: %w", err)
	}
	return nil
}

// Обновляет время жизни сессии
func (r *Redis) RefreshSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	key := sessionKeyPrefix + sessionID
	exists, err := r.conn.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("ошибка при проверке сессии: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("сессия не найдена")
	}
	success, err := r.conn.Expire(ctx, key, ttl).Result()
	if err != nil {
		return fmt.Errorf("ошибка при обновлении времени жизни сессии: %w", err)
	}
	if !success {
		// Это может произойти, если ключ был удален между Exists и Expire,
		// или если ключ не имел TTL (хотя в нашем случае он всегда с TTL).
		// Можно также проверить, существует ли ключ снова после Expire, если нужно быть абсолютно уверенным.
		return fmt.Errorf("не удалось обновить время жизни сессии (возможно, ключ был удален)")
	}
	return nil
}

// // Создает код подтверждения для email и сохраняет его в Redis
// func (r *Redis) CreateEmailCode(ctx context.Context, email string, ttl time.Duration) (string, error) {
// 	key := emailCodeKeyPrefix + email

// 	// Удаляем предыдущий код, если он существует, чтобы избежать конфликтов
// 	// и чтобы пользователь мог запросить новый код.
// 	// Ошибку удаления можно игнорировать, если ключ не существует.
// 	r.conn.Del(ctx, key) // Можно проверить err, но для этой логики не критично

// 	// Генерируем 4-значный код
// 	max := big.NewInt(10000) // 0-9999
// 	codeInt, err := rand.Int(rand.Reader, max)
// 	if err != nil {
// 		return "", fmt.Errorf("не удалось сгенерировать случайное число для кода: %w", err)
// 	}
// 	// Форматируем код до 4 цифр с ведущими нулями, если необходимо
// 	strCode := fmt.Sprintf("%04d", codeInt.Int64())

// 	// Сохраняем строковое представление кода
// 	err = r.conn.Set(ctx, key, strCode, ttl).Err()
// 	if err != nil {
// 		return "", fmt.Errorf("не удалось создать код для почты в redis: %w", err)
// 	}

// 	return strCode, nil
// }

// // Получает код подтверждения для email из Redis (в основном для внутреннего использования или тестов)
// // Для проверки кода пользователем лучше использовать VerifyEmailCode
// func (r *Redis) GetEmailCode(ctx context.Context, email string) (string, error) {
// 	key := emailCodeKeyPrefix + email
// 	code, err := r.conn.Get(ctx, key).Result()
// 	if err == redis.Nil {
// 		return "", fmt.Errorf("код для почты %s не найден или истек", email)
// 	}
// 	if err != nil {
// 		return "", fmt.Errorf("ошибка при получении кода для почты %s: %w", email, err)
// 	}
// 	return code, nil
// }

// // Удаляет код подтверждения для email из Redis
// // Может быть полезно, если пользователь отменил операцию или для очистки.
// func (r *Redis) DeleteEmailCode(ctx context.Context, email string) error {
// 	key := emailCodeKeyPrefix + email
// 	_, err := r.conn.Del(ctx, key).Result()
// 	if err == redis.Nil { // Ключа нет, значит уже удален или не существовал
// 		return nil
// 	}
// 	if err != nil {
// 		return fmt.Errorf("ошибка при удалении кода для почты %s: %w", email, err)
// 	}
// 	return nil
// }

// // --- Новые функции для верификации ---

// // VerifyEmailCode проверяет предоставленный пользователем код.
// // Если код верный, он удаляется из Redis.
// func (r *Redis) VerifyEmailCode(ctx context.Context, email string, userProvidedCode string) (bool, error) {
// 	key := emailCodeKeyPrefix + email
// 	storedCode, err := r.conn.Get(ctx, key).Result()
// 	if err == redis.Nil {
// 		return false, fmt.Errorf("код подтверждения для почты %s не найден (возможно, истек или уже использован)", email)
// 	}
// 	if err != nil {
// 		return false, fmt.Errorf("ошибка при получении кода из Redis для почты %s: %w", email, err)
// 	}

// 	if storedCode == userProvidedCode {
// 		// Код верный, удаляем его, чтобы предотвратить повторное использование
// 		_, delErr := r.conn.Del(ctx, key).Result()
// 		if delErr != nil {
// 			// Логируем ошибку, но считаем верификацию успешной, т.к. код совпал
// 			// Важно, чтобы эта ошибка не помешала пользователю
// 			fmt.Printf("Внимание: не удалось удалить использованный код для почты %s: %v\n", email, delErr)
// 			// В продакшене здесь должен быть нормальный логгер
// 		}
// 		return true, nil
// 	}

// 	return false, fmt.Errorf("неверный код подтверждения для почты %s", email)
// }

// // MarkEmailAsVerified помечает email как подтвержденный.
// // TTL здесь можно установить очень большим или 0 (без истечения срока),
// // в зависимости от вашей политики. Если 0, то ключ будет жить, пока Redis не очистится
// // или не будет удален вручную.
// func (r *Redis) MarkEmailAsVerified(ctx context.Context, email string, ttl time.Duration) error {
// 	key := verifiedEmailKeyPrefix + email
// 	// Значение может быть любым, например "1" или "true"
// 	err := r.conn.Set(ctx, key, "1", ttl).Err()
// 	if err != nil {
// 		return fmt.Errorf("не удалось пометить почту %s как подтвержденную: %w", email, err)
// 	}
// 	return nil
// }

// // IsEmailVerified проверяет, была ли почта подтверждена.
// func (r *Redis) IsEmailVerified(ctx context.Context, email string) (bool, error) {
// 	key := verifiedEmailKeyPrefix + email
// 	exists, err := r.conn.Exists(ctx, key).Result()
// 	if err != nil {
// 		return false, fmt.Errorf("ошибка при проверке статуса подтверждения почты %s: %w", email, err)
// 	}
// 	return exists == 1, nil
// }

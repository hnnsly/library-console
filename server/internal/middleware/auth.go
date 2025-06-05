package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/repository"
	"github.com/rs/zerolog/log"
)

// NewAuthMiddleware возвращает Fiber-middleware, который:
//  1. Достаёт session-ID из cookie либо "Authorization: Bearer ...".
//  2. Проверяет сессию в Redis.
//  3. Кладёт userID в c.Locals("userID").
//  4. Продлевает TTL (best-effort — ошибка логируется, но не роняет запрос).
func NewAuthMiddleware(repo *repository.LibraryRepository, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Получаем session ID
		sid := c.Cookies("session_token")
		if sid == "" {
			if ah := c.Get("Authorization"); strings.HasPrefix(ah, "Bearer ") {
				sid = strings.TrimPrefix(ah, "Bearer ")
			}
		}
		if sid == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized: empty session",
			})
		}

		// 2. Читаем сессию из Redis
		session, err := repo.GetSession(c.Context(), sid)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized: " + err.Error(),
			})
		}

		// 3. Кладём userID и role в context
		c.Locals("userID", session.UserID.String())
		c.Locals("userRole", string(session.Role))

		// 4. best-effort обновляем TTL, чтобы «подскакивала» каждая активность
		if err := repo.RefreshSession(c.Context(), sid, ttl); err != nil {
			log.Warn().Err(err).Msg("cannot refresh session TTL")
		}

		return c.Next()
	}
}

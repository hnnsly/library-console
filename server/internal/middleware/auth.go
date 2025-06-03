package middleware

import (
	"time"

	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/auth"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/rs/zerolog/log"
)

func AuthRequired(sessionManager *auth.SessionManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return httperr.New(fiber.StatusUnauthorized, "Authentication required.")
		}

		session, err := sessionManager.GetSession(c.Context(), sessionID)
		if err != nil {
			// Удаляем некорректный cookie
			c.Cookie(&fiber.Cookie{
				Name:     "session_id",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HTTPOnly: true,
				Path:     "/",
			})
			return httperr.New(fiber.StatusUnauthorized, "Invalid session.")
		}

		// Дополнительные проверки безопасности
		currentIP := c.IP()
		currentUserAgent := c.Get("User-Agent")

		if err := sessionManager.ValidateSessionSecurity(c.Context(), sessionID, currentIP, currentUserAgent); err != nil {
			log.Warn().
				Str("sessionID", sessionID).
				Str("currentIP", currentIP).
				Str("sessionIP", session.IPAddress).
				Msg("Session security validation failed")
			// В зависимости от политики безопасности можно:
			// 1. Завершить сессию
			// 2. Потребовать повторную аутентификацию
			// 3. Просто залогировать
		}

		// Обновляем время последней активности (не чаще раза в минуту)
		if time.Since(session.LastSeen) > time.Minute {
			sessionManager.UpdateLastSeen(c.Context(), sessionID)
		}

		// Сохраняем данные пользователя в контексте
		c.Locals("user_id", session.UserID)
		c.Locals("username", session.Username)
		c.Locals("role", session.Role)
		c.Locals("full_name", session.FullName)
		c.Locals("session_id", sessionID)

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return httperr.New(fiber.StatusUnauthorized, "Role information not found.")
		}

		if slices.Contains(roles, userRole) {
			return c.Next()
		}

		log.Warn().
			Str("userRole", userRole).
			Strs("requiredRoles", roles).
			Str("path", c.Path()).
			Msg("Access denied - insufficient permissions")

		return httperr.New(fiber.StatusForbidden, "Insufficient permissions.")
	}
}

// Rate limiting middleware
func RateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Здесь должна быть реализация rate limiting
		// Можно использовать готовые решения или Redis
		return c.Next()
	}
}

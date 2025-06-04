package middleware

import (
	"strings"
	"time"

	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/auth"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/rs/zerolog/log"
)

// AuthRequired проверяет аутентификацию пользователя
func AuthRequired(sessionManager *auth.SessionManager, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем session ID из cookie или Authorization header
		sessionID := getSessionID(c)
		if sessionID == "" {
			return httperr.New(fiber.StatusUnauthorized, "Authentication required.")
		}

		// Получаем данные сессии из Redis
		session, err := sessionManager.GetSession(c.Context(), sessionID)
		if err != nil {
			clearSessionCookie(c)
			return httperr.New(fiber.StatusUnauthorized, "Invalid session.")
		}

		// Сохраняем данные в контексте
		c.Locals("user_id", session.UserID)
		c.Locals("role", session.Role)
		c.Locals("session_id", sessionID)

		// Best-effort обновление TTL (не блокируем запрос при ошибке)
		if err := sessionManager.RefreshSession(c.Context(), sessionID, ttl); err != nil {
			log.Warn().Err(err).Str("session_id", sessionID).Msg("Failed to refresh session TTL")
		}

		return c.Next()
	}
}

// RequireRole проверяет, что у пользователя есть необходимая роль
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return httperr.New(fiber.StatusUnauthorized, "Role information not found.")
		}

		if !slices.Contains(roles, userRole) {
			log.Warn().
				Str("userRole", userRole).
				Strs("requiredRoles", roles).
				Str("path", c.Path()).
				Int64("userID", c.Locals("user_id").(int64)).
				Msg("Access denied - insufficient permissions")

			return httperr.New(fiber.StatusForbidden, "Insufficient permissions.")
		}

		return c.Next()
	}
}

// getSessionID извлекает session ID из cookie или Authorization header
func getSessionID(c *fiber.Ctx) string {
	// Сначала пробуем cookie
	if sessionID := c.Cookies("library-console_session_token"); sessionID != "" {
		return sessionID
	}

	// Потом Authorization header
	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	return ""
}

func clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "library-console_session_token",
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   true, // Всегда true для HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(-time.Hour),
	})
}

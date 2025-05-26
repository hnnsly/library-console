package middleware

import (
	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error" // Обновите путь
)

// RequireAuthHeader - это пример middleware для проверки заголовка авторизации.
// В реальном приложении здесь будет проверка токена (JWT, API key и т.д.).
func RequireAuthHeader() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// В реальном приложении здесь может быть более сложная логика
			// Например, проверка формата "Bearer <token>"
			return httperr.New(fiber.StatusUnauthorized, "Authorization header is missing.")
		}
		// Здесь должна быть логика валидации токена
		// if !isValidToken(authHeader) {
		//    return httperr.New(fiber.StatusUnauthorized, "Invalid authorization token.")
		// }
		return c.Next()
	}
}

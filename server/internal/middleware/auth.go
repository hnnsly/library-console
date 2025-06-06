package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/repository"
	"github.com/rs/zerolog/log"
)

func NewAuthMiddleware(repo *repository.LibraryRepository, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {

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

		session, err := repo.GetSession(c.Context(), sid)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized: " + err.Error(),
			})
		}

		c.Locals("userID", session.UserID.String())
		c.Locals("userRole", string(session.Role))

		if err := repo.RefreshSession(c.Context(), sid, ttl); err != nil {
			log.Warn().Err(err).Msg("cannot refresh session TTL")
		}

		return c.Next()
	}
}

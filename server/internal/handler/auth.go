package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User UserResponse `json:"user"`
}

func (h *Handler) login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Get user by username
	user, err := h.repo.GetUserByUsername(c.Context(), req.Username)
	if err != nil {
		log.Warn().Str("username", req.Username).Msg("Login attempt with non-existent username")
		return httperr.New(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Check if user is active
	if user.IsActive != nil && !*user.IsActive {
		return httperr.New(fiber.StatusUnauthorized, "User account is deactivated")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Warn().Str("username", req.Username).Msg("Login attempt with invalid password")
		return httperr.New(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Create session
	sessionToken, err := h.repo.CreateSession(c.Context(), user.ID, user.Role, 24*time.Hour)
	if err != nil {
		log.Error().Err(err).Str("userID", user.ID.String()).Msg("Failed to create session")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create session")
	}

	// Set session cookie
	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.Session.CookiePath,
		Value:    sessionToken,
		Expires:  time.Now().Add(h.cfg.Session.TTL),
		HTTPOnly: h.cfg.Session.HttpOnly,
		Secure:   h.cfg.Session.Secure, // Set to true in production
		SameSite: h.cfg.Session.SameSite,
	})

	response := LoginResponse{
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role),
			IsActive: user.IsActive,
		},
	}

	log.Info().Str("username", req.Username).Str("userID", user.ID.String()).Msg("User logged in successfully")

	return c.JSON(response)
}

func (h *Handler) logout(c *fiber.Ctx) error {
	// Get session token
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		if ah := c.Get("Authorization"); ah != "" && len(ah) > 7 && ah[:7] == "Bearer " {
			sessionToken = ah[7:]
		}
	}

	if sessionToken != "" {
		// Delete session from Redis
		err := h.repo.DeleteSession(c.Context(), sessionToken)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to delete session during logout")
		}
	}

	// Clear session cookie
	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.Session.CookiePath,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: h.cfg.Session.HttpOnly,
		Secure:   h.cfg.Session.Secure,
		SameSite: h.cfg.Session.SameSite,
	})

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

func (h *Handler) me(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return httperr.New(fiber.StatusUnauthorized, "User ID not found in context")
	}

	_, ok = c.Locals("userRole").(string)
	if !ok {
		return httperr.New(fiber.StatusUnauthorized, "User role not found in context")
	}

	// Parse user ID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusInternalServerError, "Invalid user ID format")
	}

	// Get user details
	user, err := h.repo.GetUserById(c.Context(), userID)
	if err != nil {
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get user for /me endpoint")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve user data")
	}

	var createdAt *string
	if user.CreatedAt != nil {
		formatted := user.CreatedAt.Format("2006-01-02T15:04:05Z")
		createdAt = &formatted
	}

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: createdAt,
	}

	return c.JSON(response)
}

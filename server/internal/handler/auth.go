package handler

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/auth"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo           *postgres.Queries
	sessionManager *auth.SessionManager
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		FullName string `json:"full_name"`
	} `json:"user"`
	SessionID string `json:"session_id"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func NewAuthHandler(repo *postgres.Queries, sessionManager *auth.SessionManager) *AuthHandler {
	return &AuthHandler{
		repo:           repo,
		sessionManager: sessionManager,
	}
}

// Валидация входных данных
func (req *LoginRequest) Validate() error {
	if strings.TrimSpace(req.Username) == "" {
		return httperr.New(fiber.StatusBadRequest, "Username is required.")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return httperr.New(fiber.StatusBadRequest, "Username must be between 3 and 50 characters.")
	}
	if req.Password == "" {
		return httperr.New(fiber.StatusBadRequest, "Password is required.")
	}
	if len(req.Password) < 6 {
		return httperr.New(fiber.StatusBadRequest, "Password is too short.")
	}
	return nil
}

func (req *ChangePasswordRequest) Validate() error {
	if req.CurrentPassword == "" {
		return httperr.New(fiber.StatusBadRequest, "Current password is required.")
	}
	if req.NewPassword == "" {
		return httperr.New(fiber.StatusBadRequest, "New password is required.")
	}
	if len(req.NewPassword) < 8 {
		return httperr.New(fiber.StatusBadRequest, "New password must be at least 8 characters long.")
	}
	// Добавляем проверку сложности пароля
	if !isPasswordStrong(req.NewPassword) {
		return httperr.New(fiber.StatusBadRequest, "New password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character.")
	}
	if req.CurrentPassword == req.NewPassword {
		return httperr.New(fiber.StatusBadRequest, "New password must be different from current password.")
	}
	return nil
}

// login авторизует пользователя
func (h *AuthHandler) login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// Валидация входных данных
	if err := req.Validate(); err != nil {
		return err
	}

	// Rate limiting проверка (можно добавить через middleware)
	clientIP := c.IP()
	if err := h.checkRateLimit(c.Context(), clientIP, req.Username); err != nil {
		return err
	}

	// Получаем пользователя из базы
	user, err := h.repo.GetUserByUsername(c.Context(), req.Username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// Логируем попытку входа с несуществующим пользователем
			log.Warn().
				Str("username", req.Username).
				Str("ip", clientIP).
				Msg("Login attempt with non-existent user")
			return httperr.New(fiber.StatusUnauthorized, "Invalid credentials.")
		}
		log.Error().Err(err).Str("username", req.Username).Msg("Failed to get user")
		return httperr.New(fiber.StatusInternalServerError, "Authentication failed.")
	}

	// Проверяем, активен ли пользователь
	if !*user.IsActive {
		log.Warn().
			Str("username", req.Username).
			Str("ip", clientIP).
			Msg("Login attempt with inactive user")
		return httperr.New(fiber.StatusUnauthorized, "Account is inactive.")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Логируем неудачную попытку входа
		log.Warn().
			Str("username", req.Username).
			Str("ip", clientIP).
			Msg("Failed login attempt")

		// Увеличиваем счетчик неудачных попыток
		h.incrementFailedAttempts(c.Context(), clientIP, req.Username)

		return httperr.New(fiber.StatusUnauthorized, "Invalid credentials.")
	}

	// Инвалидируем все существующие сессии пользователя (опционально)
	if err := h.sessionManager.InvalidateUserSessions(c.Context(), user.ID); err != nil {
		log.Error().Err(err).Int64("userID", user.ID).Msg("Failed to invalidate existing sessions")
		// Продолжаем, это не критическая ошибка
	}

	// Создаем сессию
	sessionData := auth.SessionData{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		FullName:  user.FullName,
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
		IPAddress: clientIP,
		UserAgent: c.Get("User-Agent"),
	}

	sessionID, err := h.sessionManager.CreateSession(c.Context(), sessionData)
	if err != nil {
		log.Error().Err(err).Int64("userID", user.ID).Msg("Failed to create session")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create session.")
	}

	// Обновляем время последнего входа
	err = h.repo.UpdateLastLogin(c.Context(), user.ID)
	if err != nil {
		log.Error().Err(err).Int64("userID", user.ID).Msg("Failed to update last login")
	}

	// Сбрасываем счетчик неудачных попыток
	h.resetFailedAttempts(c.Context(), clientIP, req.Username)

	// Устанавливаем cookie с улучшенной безопасностью
	h.setSecureCookie(c, sessionID)

	// Логируем успешный вход
	log.Info().
		Str("username", user.Username).
		Str("ip", clientIP).
		Str("sessionID", sessionID).
		Msg("Successful login")

	response := LoginResponse{
		SessionID: sessionID,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.Email = user.Email
	response.User.Role = user.Role
	response.User.FullName = user.FullName

	return c.JSON(response)
}

// logout выходит из системы
func (h *AuthHandler) logout(c *fiber.Ctx) error {
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		err := h.sessionManager.DeleteSession(c.Context(), sessionID)
		if err != nil {
			log.Error().Err(err).Str("sessionID", sessionID).Msg("Failed to delete session")
		}

		// Логируем выход
		if userID := c.Locals("user_id"); userID != nil {
			log.Info().
				Int64("userID", userID.(int64)).
				Str("sessionID", sessionID).
				Msg("User logged out")
		}
	}

	// Удаляем cookie
	h.clearSecureCookie(c)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// getCurrentUser получает информацию о текущем пользователе
func (h *AuthHandler) getCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	user, err := h.repo.GetUserByID(c.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		log.Error().Err(err).Int64("userID", userID).Msg("Failed to get current user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get user information.")
	}

	// Обновляем время последней активности в сессии
	sessionID := c.Cookies("session_id")
	if sessionID != "" {
		h.sessionManager.UpdateLastSeen(c.Context(), sessionID)
	}

	return c.JSON(user)
}

// changePassword изменяет пароль пользователя
func (h *AuthHandler) changePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	username := c.Locals("username").(string)

	req := new(ChangePasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// Валидация
	if err := req.Validate(); err != nil {
		return err
	}

	// Получаем текущего пользователя для проверки пароля
	user, err := h.repo.GetUserByUsername(c.Context(), username)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("Failed to get user for password change")
		return httperr.New(fiber.StatusInternalServerError, "Failed to change password.")
	}

	// Проверяем текущий пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword))
	if err != nil {
		log.Warn().
			Int64("userID", userID).
			Str("ip", c.IP()).
			Msg("Failed password change attempt - wrong current password")
		return httperr.New(fiber.StatusBadRequest, "Current password is incorrect.")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash new password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to change password.")
	}

	params := postgres.UpdatePasswordParams{
		ID:           userID,
		PasswordHash: string(hashedPassword),
	}

	err = h.repo.UpdatePassword(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("Failed to update password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to change password.")
	}

	// Инвалидируем все сессии пользователя кроме текущей
	currentSessionID := c.Cookies("session_id")
	if err := h.sessionManager.InvalidateUserSessionsExcept(c.Context(), userID, currentSessionID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("Failed to invalidate other sessions after password change")
	}

	log.Info().
		Int64("userID", userID).
		Str("ip", c.IP()).
		Msg("Password changed successfully")

	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

// Вспомогательные функы
func (h *AuthHandler) setSecureCookie(c *fiber.Ctx, sessionID string) {
	// Определяем, работаем ли мы в dev режиме
	//isProduction := c.Get("X-Forwarded-Proto") == "https" || c.Protocol() == "https"

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,  // Recommended for production
		SameSite: "Lax", // Or "Strict"
		Expires:  time.Now().Add(12 * time.Hour),
	})
}

func (h *AuthHandler) clearSecureCookie(c *fiber.Ctx) {
	isProduction := c.Get("X-Forwarded-Proto") == "https" || c.Protocol() == "https"

	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   isProduction,
		SameSite: "Lax",
		Path:     "/",
		Domain:   "",
		MaxAge:   -1, // Добавляем это для принудительного удаления
	})
}

func isPasswordStrong(password string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Rate limiting функции (упрощенная реализация)
func (h *AuthHandler) checkRateLimit(ctx context.Context, ip, username string) error {
	// Здесь должна быть реализация rate limiting через Redis
	// Пример: максимум 5 попыток за 15 минут
	return nil
}

func (h *AuthHandler) incrementFailedAttempts(ctx context.Context, ip, username string) {
	// Увеличиваем счетчик неудачных попыток в Redis
}

func (h *AuthHandler) resetFailedAttempts(ctx context.Context, ip, username string) {
	// Сбрасываем счетчик неудачных попыток в Redis
}

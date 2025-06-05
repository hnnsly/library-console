package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// Request/Response structures
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Username string            `json:"username" validate:"required,min=3,max=50"`
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=6"`
	Role     postgres.UserRole `json:"role" validate:"required"`
}

type LoginResponse struct {
	SessionToken string        `json:"session_token"`
	User         *UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID         `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Role      postgres.UserRole `json:"role"`
	IsActive  *bool             `json:"is_active"`
	CreatedAt *time.Time        `json:"created_at"`
	UpdatedAt *time.Time        `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string            `json:"username" validate:"required,min=3,max=50"`
	Email    string            `json:"email" validate:"required,email"`
	Password string            `json:"password" validate:"required,min=6"`
	Role     postgres.UserRole `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Username *string            `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string            `json:"email,omitempty" validate:"omitempty,email"`
	Role     *postgres.UserRole `json:"role,omitempty"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// Login handler
func (h *Handler) login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Get user by email
	user, err := h.repo.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusUnauthorized, "Invalid credentials.")
		}
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to get user by email")
		return httperr.New(fiber.StatusInternalServerError, "Failed to authenticate user.")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return httperr.New(fiber.StatusUnauthorized, "Invalid credentials.")
	}

	// Create session
	sessionToken, err := h.repo.CreateSession(c.Context(), user.ID, user.Role, 24*time.Hour)
	if err != nil {
		log.Error().Err(err).Str("userID", user.ID.String()).Msg("Failed to create session")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create session.")
	}

	userResp := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(&LoginResponse{
		SessionToken: sessionToken,
		User:         userResp,
	})
}

// Register handler
func (h *Handler) register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to process password.")
	}

	// Create user
	user, err := h.repo.CreateUser(c.Context(), postgres.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return httperr.New(fiber.StatusConflict, "Username or email already exists.")
		}
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create user.")
	}

	userResp := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(userResp)
}

// Logout handler
func (h *Handler) logout(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		if ah := c.Get("Authorization"); strings.HasPrefix(ah, "Bearer ") {
			sessionToken = strings.TrimPrefix(ah, "Bearer ")
		}
	}

	if sessionToken != "" {
		if err := h.repo.DeleteSession(c.Context(), sessionToken); err != nil {
			log.Warn().Err(err).Msg("Failed to delete session during logout")
		}
	}

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// Get current user handler
func (h *Handler) getCurrentUser(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return httperr.New(fiber.StatusUnauthorized, "User not authenticated.")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID.")
	}

	user, err := h.repo.GetUserByID(c.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		log.Error().Err(err).Str("userID", userID.String()).Msg("Failed to get current user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get user.")
	}

	userResp := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(userResp)
}

// List users handler
func (h *Handler) listUsers(c *fiber.Ctx) error {
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	users, err := h.repo.ListUsers(c.Context(), postgres.ListUsersParams{
		OffsetVal: int32(offset),
		LimitVal:  int32(limit),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to list users")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get users.")
	}

	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"users":  userResponses,
		"limit":  limit,
		"offset": offset,
	})
}

// Get user by ID handler
func (h *Handler) getUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID format.")
	}

	user, err := h.repo.GetUserByID(c.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to get user by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get user.")
	}

	userResp := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(userResp)
}

// Create user handler
func (h *Handler) createUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to process password.")
	}

	// Create user
	user, err := h.repo.CreateUser(c.Context(), postgres.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return httperr.New(fiber.StatusConflict, "Username or email already exists.")
		}
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create user.")
	}

	userResp := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(userResp)
}

// Update user handler
func (h *Handler) updateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID format.")
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Get current user data to fill in unchanged fields
	currentUser, err := h.repo.GetUserByID(c.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to get user for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get user.")
	}

	updateParams := postgres.UpdateUserParams{
		UserID:   userID,
		Username: currentUser.Username,
		Email:    currentUser.Email,
		Role:     currentUser.Role,
	}

	if req.Username != nil {
		updateParams.Username = *req.Username
	}
	if req.Email != nil {
		updateParams.Email = *req.Email
	}
	if req.Role != nil {
		updateParams.Role = *req.Role
	}

	user, err := h.repo.UpdateUser(c.Context(), updateParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return httperr.New(fiber.StatusConflict, "Username or email already exists.")
		}
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to update user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update user.")
	}

	userResp := &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(userResp)
}

// Deactivate user handler
func (h *Handler) deactivateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID format.")
	}

	err = h.repo.DeactivateUser(c.Context(), userID)
	if err != nil {
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to deactivate user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to deactivate user.")
	}

	return c.JSON(fiber.Map{"message": "User deactivated successfully"})
}

// Update user password handler
func (h *Handler) updateUserPassword(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID format.")
	}

	var req UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Get current user to verify old password
	user, err := h.repo.GetUserByID(c.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to get user for password update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get user.")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return httperr.New(fiber.StatusUnauthorized, "Invalid old password.")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash new password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to process new password.")
	}

	// Update password
	err = h.repo.UpdateUserPassword(c.Context(), postgres.UpdateUserPasswordParams{
		PasswordHash: string(hashedPassword),
		UserID:       userID,
	})
	if err != nil {
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to update user password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update password.")
	}

	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}

// Get users by role handler
func (h *Handler) getUsersByRole(c *fiber.Ctx) error {
	roleStr := c.Params("role")
	role := postgres.UserRole(roleStr)

	// Validate role
	if role != postgres.UserRoleAdministrator &&
		role != postgres.UserRoleLibrarian &&
		role != postgres.UserRoleReader {
		return httperr.New(fiber.StatusBadRequest, "Invalid role. Must be one of: administrator, librarian, reader")
	}

	users, err := h.repo.ListUsersByRole(c.Context(), role)
	if err != nil {
		log.Error().Err(err).Str("role", roleStr).Msg("Failed to get users by role")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get users.")
	}

	userResponses := make([]*UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return c.JSON(fiber.Map{
		"users": userResponses,
		"role":  role,
		"count": len(userResponses),
	})
}

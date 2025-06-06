package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  *bool     `json:"is_active"`
	CreatedAt *string   `json:"created_at"`
}

func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	users, err := h.repo.GetAllUsers(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve users")
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		var createdAt *string
		if user.CreatedAt != nil {
			formatted := user.CreatedAt.Format("2006-01-02T15:04:05Z")
			createdAt = &formatted
		}

		response[i] = UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: createdAt,
		}
	}

	return c.JSON(response)
}

func (h *Handler) getUserById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID format")
	}

	user, err := h.repo.GetUserById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found")
		}
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to get user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve user")
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

func (h *Handler) createUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate role
	validRoles := map[string]bool{
		"administrator": true,
		"librarian":     true,
	}
	if !validRoles[req.Role] {
		return httperr.New(fiber.StatusBadRequest, "Invalid role value")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to process password")
	}

	user, err := h.repo.CreateUser(c.Context(), postgres.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         postgres.UserRole(req.Role),
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return httperr.New(fiber.StatusConflict, "User with this username or email already exists")
		}
		log.Error().Err(err).Msg("Failed to create user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create user")
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

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *Handler) updateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID format")
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate role
	validRoles := map[string]bool{
		"administrator": true,
		"librarian":     true,
	}
	if !validRoles[req.Role] {
		return httperr.New(fiber.StatusBadRequest, "Invalid role value")
	}

	user, err := h.repo.UpdateUser(c.Context(), postgres.UpdateUserParams{
		ID:    id,
		Email: req.Email,
		Role:  postgres.UserRole(req.Role),
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found")
		}
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to update user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update user")
	}

	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
		IsActive: user.IsActive,
	}

	return c.JSON(response)
}

func (h *Handler) deactivateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID format")
	}

	err = h.repo.DeactivateUser(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("userID", idStr).Msg("Failed to deactivate user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to deactivate user")
	}

	return c.JSON(fiber.Map{"message": "User deactivated successfully"})
}

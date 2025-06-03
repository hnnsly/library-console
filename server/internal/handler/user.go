package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/auth"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo           *postgres.Queries
	sessionManager *auth.SessionManager
}

type CreateUserRequest struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
	FullName string  `json:"full_name"`
	Phone    *string `json:"phone"`
}

type UpdateUserRequest struct {
	Email    string  `json:"email"`
	FullName string  `json:"full_name"`
	Phone    *string `json:"phone"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

type GetAllUsersRequest struct {
	PageOffset int32 `json:"page_offset"`
	PageLimit  int32 `json:"page_limit"`
}

func NewUserHandler(repo *postgres.Queries, sessionManager *auth.SessionManager) *UserHandler {
	return &UserHandler{
		repo:           repo,
		sessionManager: sessionManager,
	}
}

// CreateUser создает нового пользователя (только для администраторов)
func (h *UserHandler) createUser(c *fiber.Ctx) error {
	req := new(CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Добавить валидацию полей

	// Только super_admin может создавать admin'ов
	currentUserRole := c.Locals("role").(string)
	if req.Role == "admin" && currentUserRole != "super_admin" {
		return httperr.New(fiber.StatusForbidden, "Only super administrators can create administrators.")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create user.")
	}

	currentUserID := c.Locals("user_id").(int64)
	params := postgres.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		FullName:     req.FullName,
		Phone:        req.Phone,
		CreatedBy:    &currentUserID,
	}

	user, err := h.repo.CreateUser(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "username") {
				return httperr.New(fiber.StatusConflict, "Username already exists.")
			}
			if strings.Contains(err.Error(), "email") {
				return httperr.New(fiber.StatusConflict, "Email already exists.")
			}
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create user.")
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetAllUsers получает список всех пользователей
func (h *UserHandler) getAllUsers(c *fiber.Ctx) error {
	req := new(GetAllUsersRequest)
	if err := c.BodyParser(req); err != nil {
		// Если тело запроса пустое или некорректное, используем значения по умолчанию
		req.PageOffset = 0
		req.PageLimit = 20
	}

	if req.PageLimit == 0 {
		req.PageLimit = 20
	}
	if req.PageLimit > 100 {
		req.PageLimit = 100
	}

	params := postgres.GetAllUsersParams{
		OffsetUsers: req.PageOffset,
		LimitUsers:  req.PageLimit,
	}

	users, err := h.repo.GetAllUsers(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve users.")
	}

	if users == nil {
		users = []*postgres.GetAllUsersRow{}
	}

	return c.JSON(users)
}

// GetUserByID получает пользователя по ID
func (h *UserHandler) getUserByID(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	user, err := h.repo.GetUserByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		log.Error().Err(err).Int64("userID", id).Msg("Failed to get user by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve user.")
	}

	return c.JSON(user)
}

// UpdateUser обновляет информацию о пользователе
func (h *UserHandler) updateUser(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Добавить валидацию полей

	// Пользователи могут редактировать только себя, администраторы - всех
	currentUserID := c.Locals("user_id").(int64)
	currentUserRole := c.Locals("role").(string)

	if currentUserID != id && currentUserRole != "admin" && currentUserRole != "super_admin" {
		return httperr.New(fiber.StatusForbidden, "You can only update your own profile.")
	}

	params := postgres.UpdateUserParams{
		ID:       id,
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
	}

	user, err := h.repo.UpdateUser(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Email already exists.")
		}
		log.Error().Err(err).Int64("userID", id).Msg("Failed to update user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update user.")
	}

	return c.JSON(user)
}

// UpdateUserRole обновляет роль пользователя (только для администраторов)
func (h *UserHandler) updateUserRole(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateUserRoleRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Валидация роли

	// Только super_admin может назначать роль admin
	currentUserRole := c.Locals("role").(string)
	if req.Role == "admin" && currentUserRole != "super_admin" {
		return httperr.New(fiber.StatusForbidden, "Only super administrators can assign admin role.")
	}

	params := postgres.UpdateUserRoleParams{
		ID:   id,
		Role: req.Role,
	}

	err = h.repo.UpdateUserRole(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "User not found or cannot modify first admin.")
		}
		log.Error().Err(err).Int64("userID", id).Msg("Failed to update user role")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update user role.")
	}

	// Инвалидируем все сессии пользователя при изменении роли
	if err := h.sessionManager.InvalidateAllUserSessions(c.Context(), id); err != nil {
		log.Error().Err(err).Int64("userID", id).Msg("Failed to invalidate user sessions after role change")
	}

	return c.JSON(fiber.Map{"message": "User role updated successfully"})
}

// DeactivateUser деактивирует пользователя
func (h *UserHandler) deactivateUser(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.DeactivateUser(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "User not found or cannot deactivate first admin.")
		}
		log.Error().Err(err).Int64("userID", id).Msg("Failed to deactivate user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to deactivate user.")
	}

	// Инвалидируем все сессии деактивированного пользователя
	if err := h.sessionManager.InvalidateAllUserSessions(c.Context(), id); err != nil {
		log.Error().Err(err).Int64("userID", id).Msg("Failed to invalidate sessions for deactivated user")
	}

	return c.JSON(fiber.Map{"message": "User deactivated successfully"})
}

// ActivateUser активирует пользователя
func (h *UserHandler) activateUser(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.ActivateUser(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "User not found.")
		}
		log.Error().Err(err).Int64("userID", id).Msg("Failed to activate user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to activate user.")
	}

	return c.JSON(fiber.Map{"message": "User activated successfully"})
}

// DeleteUser удаляет пользователя (только для super_admin)
func (h *UserHandler) deleteUser(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	// Инвалидируем все сессии пользователя перед удалением
	if err := h.sessionManager.InvalidateAllUserSessions(c.Context(), id); err != nil {
		log.Error().Err(err).Int64("userID", id).Msg("Failed to invalidate sessions before user deletion")
	}

	err = h.repo.DeleteUser(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "User not found or cannot delete first admin.")
		}
		log.Error().Err(err).Int64("userID", id).Msg("Failed to delete user")
		return httperr.New(fiber.StatusInternalServerError, "Failed to delete user.")
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// GetUsersByRole получает пользователей по роли
func (h *UserHandler) GetUsersByRole(c *fiber.Ctx) error {
	role := c.Params("role")
	if role == "" {
		return httperr.New(fiber.StatusBadRequest, "Role is required.")
	}

	// TODO: Валидация роли

	users, err := h.repo.GetUsersByRole(c.Context(), role)
	if err != nil {
		log.Error().Err(err).Str("role", role).Msg("Failed to get users by role")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve users.")
	}

	if users == nil {
		users = []*postgres.GetUsersByRoleRow{}
	}

	return c.JSON(users)
}

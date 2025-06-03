package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

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

// createUser создает нового пользователя (только для администраторов)
func (h *Handler) createUser(c *fiber.Ctx) error {
	req := new(CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: username, email, password, role, full_name
	// TODO: Validate username format (alphanumeric, 3-50 chars)
	// TODO: Validate email format
	// TODO: Validate password strength (min 8 chars, complexity)
	// TODO: Validate role is one of: admin, librarian
	// TODO: Validate full_name min length 2, max length 100
	// TODO: Validate phone format if provided

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

// getAllUsers получает список всех пользователей
func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	req := new(GetAllUsersRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	if req.PageLimit == 0 {
		req.PageLimit = 20 // default limit
	}

	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0

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

// getUserByID получает пользователя по ID
func (h *Handler) getUserByID(c *fiber.Ctx) error {
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

// updateUser обновляет информацию о пользователе
func (h *Handler) updateUser(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: email, full_name
	// TODO: Validate email format
	// TODO: Validate full_name min length 2, max length 100
	// TODO: Validate phone format if provided

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

// updateUserRole обновляет роль пользователя (только для администраторов)
func (h *Handler) updateUserRole(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateUserRoleRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate role is one of: admin, librarian

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

	return c.JSON(fiber.Map{"message": "User role updated successfully"})
}

// deactivateUser деактивирует пользователя
func (h *Handler) deactivateUser(c *fiber.Ctx) error {
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

	return c.JSON(fiber.Map{"message": "User deactivated successfully"})
}

// activateUser активирует пользователя
func (h *Handler) activateUser(c *fiber.Ctx) error {
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

// deleteUser удаляет пользователя (только для super_admin)
func (h *Handler) deleteUser(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
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

// getUsersByRole получает пользователей по роли
func (h *Handler) getUsersByRole(c *fiber.Ctx) error {
	role := c.Params("role")
	if role == "" {
		return httperr.New(fiber.StatusBadRequest, "Role is required.")
	}

	// TODO: Validate role is one of: super_admin, admin, librarian

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

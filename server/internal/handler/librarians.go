package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type CreateLibrarianRequest struct {
	FullName   string  `json:"full_name"`
	EmployeeID string  `json:"employee_id"`
	Position   string  `json:"position"`
	Phone      *string `json:"phone"`
	Email      *string `json:"email"`
}

type UpdateLibrarianRequest struct {
	FullName string  `json:"full_name"`
	Position string  `json:"position"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
}

// createLibrarian создает нового библиотекаря
func (h *Handler) createLibrarian(c *fiber.Ctx) error {
	req := new(CreateLibrarianRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: full_name, employee_id, position
	// TODO: Validate full_name min length 2, max length 100
	// TODO: Validate employee_id is unique and matches format
	// TODO: Validate position is one of: librarian, senior_librarian, head_librarian, assistant
	// TODO: Validate phone format if provided
	// TODO: Validate email format if provided

	params := postgres.CreateLibrarianParams{
		FullName:   req.FullName,
		EmployeeID: req.EmployeeID,
		Position:   req.Position,
		Phone:      req.Phone,
		Email:      req.Email,
	}

	librarian, err := h.repo.CreateLibrarian(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create librarian")
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Librarian with this employee ID already exists.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create librarian.")
	}

	return c.Status(fiber.StatusCreated).JSON(librarian)
}

// getLibrarianByID получает библиотекаря по ID
func (h *Handler) getLibrarianByID(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	librarian, err := h.repo.GetLibrarianByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Librarian not found.")
		}
		log.Error().Err(err).Int64("librarianID", id).Msg("Failed to get librarian by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve librarian.")
	}

	return c.JSON(librarian)
}

// getLibrarianByEmployeeID получает библиотекаря по табельному номеру
func (h *Handler) getLibrarianByEmployeeID(c *fiber.Ctx) error {
	employeeID := c.Params("employee_id")
	if employeeID == "" {
		return httperr.New(fiber.StatusBadRequest, "Employee ID is required.")
	}

	// TODO: Validate employee_id format

	librarian, err := h.repo.GetLibrarianByEmployeeID(c.Context(), employeeID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Librarian not found.")
		}
		log.Error().Err(err).Str("employeeID", employeeID).Msg("Failed to get librarian by employee ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve librarian.")
	}

	return c.JSON(librarian)
}

// getAllLibrarians получает список всех активных библиотекарей
func (h *Handler) getAllLibrarians(c *fiber.Ctx) error {
	librarians, err := h.repo.GetAllLibrarians(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all librarians")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve librarians.")
	}

	if librarians == nil {
		librarians = []*postgres.Librarian{}
	}

	return c.JSON(librarians)
}

// updateLibrarian обновляет информацию о библиотекаре
func (h *Handler) updateLibrarian(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateLibrarianRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: full_name, position
	// TODO: Validate full_name min length 2, max length 100
	// TODO: Validate position is one of: librarian, senior_librarian, head_librarian, assistant
	// TODO: Validate phone format if provided
	// TODO: Validate email format if provided
	// TODO: Validate librarian exists and is active

	params := postgres.UpdateLibrarianParams{
		LibrarianID: id,
		FullName:    req.FullName,
		Position:    req.Position,
		Phone:       req.Phone,
		Email:       req.Email,
	}

	librarian, err := h.repo.UpdateLibrarian(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Librarian not found.")
		}
		log.Error().Err(err).Int64("librarianID", id).Msg("Failed to update librarian")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update librarian.")
	}

	return c.JSON(librarian)
}

// deactivateLibrarian деактивирует библиотекаря (увольнение)
func (h *Handler) deactivateLibrarian(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.DeactivateLibrarian(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Librarian not found.")
		}
		log.Error().Err(err).Int64("librarianID", id).Msg("Failed to deactivate librarian")
		return httperr.New(fiber.StatusInternalServerError, "Failed to deactivate librarian.")
	}

	return c.JSON(fiber.Map{"message": "Librarian successfully deactivated"})
}

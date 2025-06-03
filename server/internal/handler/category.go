package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type CreateCategoryRequest struct {
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	DefaultLoanDays int     `json:"default_loan_days"`
}

// createCategory создает новую категорию книг
func (h *Handler) createCategory(c *fiber.Ctx) error {
	req := new(CreateCategoryRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required field: name (min length 2, max length 100)
	// TODO: Validate name is unique
	// TODO: Validate default_loan_days > 0 and <= 365
	// TODO: Validate description max length 500 if provided

	params := postgres.CreateCategoryParams{
		Name:            req.Name,
		Description:     req.Description,
		DefaultLoanDays: req.DefaultLoanDays,
	}

	category, err := h.repo.CreateCategory(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create category")
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Category with this name already exists.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create category.")
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// getAllCategories получает список всех категорий книг
func (h *Handler) getAllCategories(c *fiber.Ctx) error {
	categories, err := h.repo.GetAllCategories(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all categories")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve categories.")
	}

	if categories == nil {
		categories = []*postgres.BookCategory{}
	}

	return c.JSON(categories)
}

// getCategoryStatistics получает статистику по категориям книг
func (h *Handler) getCategoryStatistics(c *fiber.Ctx) error {
	statistics, err := h.repo.GetCategoryStatistics(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get category statistics")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve category statistics.")
	}

	if statistics == nil {
		statistics = []*postgres.GetCategoryStatisticsRow{}
	}

	return c.JSON(statistics)
}

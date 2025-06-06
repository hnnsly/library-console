package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

type CreateAuthorRequest struct {
	FullName string `json:"full_name" validate:"required"`
}

func (h *Handler) getAllAuthors(c *fiber.Ctx) error {
	authors, err := h.repo.GetAllAuthors(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all authors")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve authors")
	}

	return c.JSON(authors)
}

func (h *Handler) searchAuthors(c *fiber.Ctx) error {
	searchTerm := c.Query("q")
	if searchTerm == "" {
		return httperr.New(fiber.StatusBadRequest, "Search term is required")
	}

	authors, err := h.repo.SearchAuthors(c.Context(), &searchTerm)
	if err != nil {
		log.Error().Err(err).Str("searchTerm", searchTerm).Msg("Failed to search authors")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search authors")
	}

	return c.JSON(authors)
}

func (h *Handler) getAuthorById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format")
	}

	author, err := h.repo.GetAuthorById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Author not found")
		}
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to get author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve author")
	}

	return c.JSON(author)
}

func (h *Handler) createAuthor(c *fiber.Ctx) error {
	var req CreateAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	author, err := h.repo.CreateAuthor(c.Context(), req.FullName)
	if err != nil {
		log.Error().Err(err).Str("fullName", req.FullName).Msg("Failed to create author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create author")
	}

	return c.Status(fiber.StatusCreated).JSON(author)
}

func (h *Handler) getAuthorBooks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format")
	}

	books, err := h.repo.GetAuthorBooks(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to get author books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve author books")
	}

	return c.JSON(books)
}

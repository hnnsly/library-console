package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

type CreateBookCopyRequest struct {
	BookID       string  `json:"book_id" validate:"required"`
	CopyCode     string  `json:"copy_code" validate:"required"`
	HallID       *string `json:"hall_id"`
	LocationInfo *string `json:"location_info"`
}

type UpdateBookCopyStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

func (h *Handler) getBookCopiesByBookId(c *fiber.Ctx) error {
	bookIdStr := c.Params("bookId")
	bookId, err := uuid.Parse(bookIdStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format")
	}

	copies, err := h.repo.GetBookCopiesByBookId(c.Context(), bookId)
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIdStr).Msg("Failed to get book copies")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copies")
	}

	return c.JSON(copies)
}

func (h *Handler) getBookCopiesByHall(c *fiber.Ctx) error {
	hallIdStr := c.Params("hallId")
	var hallId *uuid.UUID

	if hallIdStr != "" {
		id, err := uuid.Parse(hallIdStr)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
		}
		hallId = &id
	}

	copies, err := h.repo.GetBookCopiesByHall(c.Context(), hallId)
	if err != nil {
		log.Error().Err(err).Str("hallID", hallIdStr).Msg("Failed to get book copies by hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copies")
	}

	return c.JSON(copies)
}

func (h *Handler) getBookCopyByCode(c *fiber.Ctx) error {
	copyCode := c.Params("copyCode")
	if copyCode == "" {
		return httperr.New(fiber.StatusBadRequest, "Copy code is required")
	}

	copy, err := h.repo.GetBookCopyByCode(c.Context(), copyCode)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found")
		}
		log.Error().Err(err).Str("copyCode", copyCode).Msg("Failed to get book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy")
	}

	return c.JSON(copy)
}

func (h *Handler) getBookCopyById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid copy ID format")
	}

	copy, err := h.repo.GetBookCopyById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found")
		}
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to get book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy")
	}

	return c.JSON(copy)
}

func (h *Handler) createBookCopy(c *fiber.Ctx) error {
	var req CreateBookCopyRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	bookId, err := uuid.Parse(req.BookID)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format")
	}

	var hallId *uuid.UUID
	if req.HallID != nil && *req.HallID != "" {
		id, err := uuid.Parse(*req.HallID)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
		}
		hallId = &id
	}

	copy, err := h.repo.CreateBookCopy(c.Context(), postgres.CreateBookCopyParams{
		BookID:       bookId,
		CopyCode:     req.CopyCode,
		HallID:       hallId,
		LocationInfo: req.LocationInfo,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create book copy")
	}

	return c.Status(fiber.StatusCreated).JSON(copy)
}

func (h *Handler) updateBookCopyStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid copy ID format")
	}

	var req UpdateBookCopyStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate status
	validStatuses := map[string]bool{
		"available": true,
		"issued":    true,
		"reserved":  true,
		"lost":      true,
		"damaged":   true,
	}
	if !validStatuses[req.Status] {
		return httperr.New(fiber.StatusBadRequest, "Invalid status value")
	}

	err = h.repo.UpdateBookCopyStatus(c.Context(), postgres.UpdateBookCopyStatusParams{
		CopyID: id,
		Status: postgres.NullBookStatus{
			BookStatus: postgres.BookStatus(req.Status),
			Valid:      true,
		},
	})
	if err != nil {
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to update book copy status")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book copy status")
	}

	return c.JSON(fiber.Map{"message": "Book copy status updated successfully"})
}

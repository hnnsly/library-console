package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

type CreateReaderRequest struct {
	TicketNumber string  `json:"ticket_number" validate:"required"`
	FullName     string  `json:"full_name" validate:"required"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
}

type UpdateReaderRequest struct {
	FullName string  `json:"full_name" validate:"required"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
}

func (h *Handler) getActiveReaders(c *fiber.Ctx) error {
	readers, err := h.repo.GetActiveReaders(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get active readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve readers")
	}

	return c.JSON(readers)
}

func (h *Handler) searchReaders(c *fiber.Ctx) error {
	searchTerm := c.Query("q")
	includeInactiveStr := c.Query("include_inactive", "false")

	includeInactive, err := strconv.ParseBool(includeInactiveStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid include_inactive parameter")
	}

	readers, err := h.repo.SearchReaders(c.Context(), postgres.SearchReadersParams{
		SearchTerm:      &searchTerm,
		IncludeInactive: includeInactive,
	})
	if err != nil {
		log.Error().Err(err).Str("searchTerm", searchTerm).Msg("Failed to search readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search readers")
	}

	return c.JSON(readers)
}

func (h *Handler) getReaderById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	reader, err := h.repo.GetReaderById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found")
		}
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to get reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader")
	}

	return c.JSON(reader)
}

func (h *Handler) getReaderByTicketNumber(c *fiber.Ctx) error {
	ticketNumber := c.Params("ticketNumber")
	if ticketNumber == "" {
		return httperr.New(fiber.StatusBadRequest, "Ticket number is required")
	}

	reader, err := h.repo.GetReaderByTicketNumber(c.Context(), ticketNumber)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found")
		}
		log.Error().Err(err).Str("ticketNumber", ticketNumber).Msg("Failed to get reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader")
	}

	return c.JSON(reader)
}

func (h *Handler) createReader(c *fiber.Ctx) error {
	var req CreateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	reader, err := h.repo.CreateReader(c.Context(), postgres.CreateReaderParams{
		TicketNumber: req.TicketNumber,
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return httperr.New(fiber.StatusConflict, "Reader with this ticket number already exists")
		}
		log.Error().Err(err).Msg("Failed to create reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create reader")
	}

	return c.Status(fiber.StatusCreated).JSON(reader)
}

func (h *Handler) updateReader(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	var req UpdateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	reader, err := h.repo.UpdateReader(c.Context(), postgres.UpdateReaderParams{
		ID:       id,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found")
		}
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to update reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update reader")
	}

	return c.JSON(reader)
}

func (h *Handler) deactivateReader(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	err = h.repo.DeactivateReader(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to deactivate reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to deactivate reader")
	}

	return c.JSON(fiber.Map{"message": "Reader deactivated successfully"})
}

func (h *Handler) getReaderActiveBooks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	books, err := h.repo.GetReaderActiveBooks(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to get reader active books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader active books")
	}

	return c.JSON(books)
}

func (h *Handler) getReaderFines(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	fines, err := h.repo.GetReaderFines(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to get reader fines")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader fines")
	}

	return c.JSON(fines)
}

func (h *Handler) getReaderVisitHistory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	visits, err := h.repo.GetReaderVisitHistory(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to get reader visit history")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader visit history")
	}

	return c.JSON(visits)
}

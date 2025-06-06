package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

type CreateFineRequest struct {
	ReaderID    string  `json:"reader_id" validate:"required"`
	BookIssueID *string `json:"book_issue_id"`
	Amount      string  `json:"amount" validate:"required"`
	Reason      string  `json:"reason" validate:"required"`
}

func (h *Handler) getUnpaidFines(c *fiber.Ctx) error {
	fines, err := h.repo.GetUnpaidFines(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get unpaid fines")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve unpaid fines")
	}

	return c.JSON(fines)
}

func (h *Handler) createFine(c *fiber.Ctx) error {
	var req CreateFineRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	readerID, err := uuid.Parse(req.ReaderID)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	var bookIssueID *uuid.UUID
	if req.BookIssueID != nil && *req.BookIssueID != "" {
		id, err := uuid.Parse(*req.BookIssueID)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid book issue ID format")
		}
		bookIssueID = &id
	}

	amount, err := decimal.Parse(req.Amount)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid amount format")
	}

	fine, err := h.repo.CreateFine(c.Context(), postgres.CreateFineParams{
		ReaderID:    readerID,
		BookIssueID: bookIssueID,
		Amount:      amount,
		Reason:      req.Reason,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create fine")
	}

	return c.Status(fiber.StatusCreated).JSON(fine)
}

func (h *Handler) payFine(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid fine ID format")
	}

	fine, err := h.repo.PayFine(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Fine not found")
		}
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to pay fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to pay fine")
	}

	return c.JSON(fine)
}

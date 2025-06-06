package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

type RegisterHallEntryRequest struct {
	TicketNumber string `json:"ticket_number" validate:"required"`
	HallID       string `json:"hall_id" validate:"required"`
}

type RegisterHallExitRequest struct {
	TicketNumber string `json:"ticket_number" validate:"required"`
	HallID       string `json:"hall_id" validate:"required"`
}

func (h *Handler) getRecentHallVisits(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "50")
	daysStr := c.Query("days", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid limit parameter")
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid days parameter")
	}

	sinceDate := time.Now().AddDate(0, 0, -days)

	visits, err := h.repo.GetRecentHallVisits(c.Context(), postgres.GetRecentHallVisitsParams{
		SinceDate:  &sinceDate,
		LimitCount: int32(limit),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent hall visits")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve recent hall visits")
	}

	return c.JSON(visits)
}

func (h *Handler) registerHallEntry(c *fiber.Ctx) error {
	var req RegisterHallEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	hallID, err := uuid.Parse(req.HallID)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
	}

	// Get librarian ID from context
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return httperr.New(fiber.StatusUnauthorized, "User ID not found in context")
	}
	librarianID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusInternalServerError, "Invalid user ID format")
	}

	entry, err := h.repo.RegisterHallEntry(c.Context(), postgres.RegisterHallEntryParams{
		TicketNumber: req.TicketNumber,
		HallID:       hallID,
		LibrarianID:  &librarianID,
	})
	if err != nil {
		log.Error().Err(err).Str("ticketNumber", req.TicketNumber).Msg("Failed to register hall entry")
		return httperr.New(fiber.StatusInternalServerError, "Failed to register hall entry")
	}

	// Update hall visitor count (increment)
	err = h.repo.UpdateHallVisitorCount(c.Context(), postgres.UpdateHallVisitorCountParams{
		Change: &[]int{1}[0],
		HallID: hallID,
	})
	if err != nil {
		log.Warn().Err(err).Str("hallID", req.HallID).Msg("Failed to update hall visitor count")
	}

	return c.Status(fiber.StatusCreated).JSON(entry)
}

func (h *Handler) registerHallExit(c *fiber.Ctx) error {
	var req RegisterHallExitRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	hallID, err := uuid.Parse(req.HallID)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
	}

	// Get librarian ID from context
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return httperr.New(fiber.StatusUnauthorized, "User ID not found in context")
	}
	librarianID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusInternalServerError, "Invalid user ID format")
	}

	exit, err := h.repo.RegisterHallExit(c.Context(), postgres.RegisterHallExitParams{
		TicketNumber: req.TicketNumber,
		HallID:       hallID,
		LibrarianID:  &librarianID,
	})
	if err != nil {
		log.Error().Err(err).Str("ticketNumber", req.TicketNumber).Msg("Failed to register hall exit")
		return httperr.New(fiber.StatusInternalServerError, "Failed to register hall exit")
	}

	// Update hall visitor count (decrement)
	err = h.repo.UpdateHallVisitorCount(c.Context(), postgres.UpdateHallVisitorCountParams{
		Change: &[]int{-1}[0],
		HallID: hallID,
	})
	if err != nil {
		log.Warn().Err(err).Str("hallID", req.HallID).Msg("Failed to update hall visitor count")
	}

	return c.Status(fiber.StatusCreated).JSON(exit)
}

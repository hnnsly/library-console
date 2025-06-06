package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

type CreateReadingHallRequest struct {
	HallName       string  `json:"hall_name" validate:"required"`
	Specialization *string `json:"specialization"`
	TotalSeats     int     `json:"total_seats" validate:"required,min=1"`
}

type UpdateReadingHallRequest struct {
	HallName       string  `json:"hall_name" validate:"required"`
	Specialization *string `json:"specialization"`
	TotalSeats     int     `json:"total_seats" validate:"required,min=1"`
}

func (h *Handler) getAllReadingHalls(c *fiber.Ctx) error {
	halls, err := h.repo.GetAllReadingHalls(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all reading halls")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reading halls")
	}

	return c.JSON(halls)
}

func (h *Handler) getHallsDashboard(c *fiber.Ctx) error {
	dashboard, err := h.repo.GetHallsDashboard(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get halls dashboard")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve halls dashboard")
	}

	return c.JSON(dashboard)
}

func (h *Handler) getReadingHallById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
	}

	hall, err := h.repo.GetReadingHallById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reading hall not found")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to get reading hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reading hall")
	}

	return c.JSON(hall)
}

func (h *Handler) createReadingHall(c *fiber.Ctx) error {
	var req CreateReadingHallRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	hall, err := h.repo.CreateReadingHall(c.Context(), postgres.CreateReadingHallParams{
		HallName:       req.HallName,
		Specialization: req.Specialization,
		TotalSeats:     req.TotalSeats,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return httperr.New(fiber.StatusConflict, "Reading hall with this name already exists")
		}
		log.Error().Err(err).Msg("Failed to create reading hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create reading hall")
	}

	return c.Status(fiber.StatusCreated).JSON(hall)
}

func (h *Handler) updateReadingHall(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
	}

	var req UpdateReadingHallRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	hall, err := h.repo.UpdateReadingHall(c.Context(), postgres.UpdateReadingHallParams{
		ID:             id,
		HallName:       req.HallName,
		Specialization: req.Specialization,
		TotalSeats:     req.TotalSeats,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reading hall not found")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to update reading hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update reading hall")
	}

	return c.JSON(hall)
}

func (h *Handler) getDailyVisitStats(c *fiber.Ctx) error {
	hallIdStr := c.Params("id")
	hallId, err := uuid.Parse(hallIdStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid start_date format, use YYYY-MM-DD")
		}
		startDate = &parsed
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid end_date format, use YYYY-MM-DD")
		}
		endDate = &parsed
	}

	stats, err := h.repo.GetDailyVisitStats(c.Context(), postgres.GetDailyVisitStatsParams{
		HallID:    hallId,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		log.Error().Err(err).Str("hallID", hallIdStr).Msg("Failed to get daily visit stats")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve daily visit stats")
	}

	return c.JSON(stats)
}

func (h *Handler) getHourlyVisitStats(c *fiber.Ctx) error {
	hallIdStr := c.Params("id")
	hallId, err := uuid.Parse(hallIdStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format")
	}

	visitDateStr := c.Query("date")
	var visitDate *time.Time
	if visitDateStr != "" {
		parsed, err := time.Parse("2006-01-02", visitDateStr)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid date format, use YYYY-MM-DD")
		}
		visitDate = &parsed
	} else {
		now := time.Now()
		visitDate = &now
	}

	stats, err := h.repo.GetHourlyVisitStats(c.Context(), postgres.GetHourlyVisitStatsParams{
		HallID:    hallId,
		VisitDate: visitDate,
	})
	if err != nil {
		log.Error().Err(err).Str("hallID", hallIdStr).Msg("Failed to get hourly visit stats")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve hourly visit stats")
	}

	return c.JSON(stats)
}

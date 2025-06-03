package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

// getAllHalls получает список всех открытых читальных залов с информацией о занятости
func (h *Handler) getAllHalls(c *fiber.Ctx) error {
	halls, err := h.repo.GetAllHalls(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all halls")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve halls.")
	}

	if halls == nil {
		halls = []*postgres.GetAllHallsRow{}
	}

	return c.JSON(halls)
}

// getHallByID получает детальную информацию о читальном зале по ID
func (h *Handler) getHallByID(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	hall, err := h.repo.GetHallByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Hall not found.")
		}
		log.Error().Err(err).Int64("hallID", id).Msg("Failed to get hall by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve hall.")
	}

	return c.JSON(hall)
}

// getHallStatistics получает статистику по всем читальным залам
func (h *Handler) getHallStatistics(c *fiber.Ctx) error {
	statistics, err := h.repo.GetHallStatistics(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get hall statistics")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve hall statistics.")
	}

	if statistics == nil {
		statistics = []*postgres.GetHallStatisticsRow{}
	}

	return c.JSON(statistics)
}

// updateHallOccupancy обновляет информацию о занятости читального зала
func (h *Handler) updateHallOccupancy(c *fiber.Ctx) error {
	hallIDStr := c.Params("id")
	if hallIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Hall ID is required.")
	}

	hallID, err := strconv.Atoi(hallIDStr)
	if err != nil || hallID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID.")
	}

	// TODO: Validate hall_id exists

	err = h.repo.UpdateHallOccupancy(c.Context(), hallID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Hall not found.")
		}
		log.Error().Err(err).Int("hallID", hallID).Msg("Failed to update hall occupancy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update hall occupancy.")
	}

	return c.JSON(fiber.Map{"message": "Hall occupancy updated successfully"})
}

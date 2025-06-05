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

// Request/Response structs для reading halls
type CreateReadingHallRequest struct {
	LibraryName    string  `json:"library_name" validate:"required,min=2,max=100"`
	HallName       string  `json:"hall_name" validate:"required,min=2,max=100"`
	Specialization *string `json:"specialization" validate:"omitempty,max=200"`
	TotalSeats     int     `json:"total_seats" validate:"required,min=1,max=10000"`
}

type UpdateReadingHallRequest struct {
	LibraryName    *string `json:"library_name" validate:"omitempty,min=2,max=100"`
	HallName       *string `json:"hall_name" validate:"omitempty,min=2,max=100"`
	Specialization *string `json:"specialization" validate:"omitempty,max=200"`
	TotalSeats     *int    `json:"total_seats" validate:"omitempty,min=1,max=10000"`
	OccupiedSeats  *int    `json:"occupied_seats" validate:"omitempty,min=0"`
}

type UpdateHallOccupancyRequest struct {
	OccupiedSeats *int `json:"occupied_seats" validate:"omitempty,min=0"`
}

type ReadingHallResponse struct {
	ID             uuid.UUID  `json:"id"`
	LibraryName    string     `json:"library_name"`
	HallName       string     `json:"hall_name"`
	Specialization *string    `json:"specialization"`
	TotalSeats     int        `json:"total_seats"`
	OccupiedSeats  *int       `json:"occupied_seats"`
	FreeSeats      *int       `json:"free_seats,omitempty"`
	CreatedAt      *time.Time `json:"created_at"`
}

type HallStatisticsResponse struct {
	ID                uuid.UUID `json:"id"`
	LibraryName       string    `json:"library_name"`
	HallName          string    `json:"hall_name"`
	Specialization    *string   `json:"specialization"`
	TotalSeats        int       `json:"total_seats"`
	OccupiedSeats     *int      `json:"occupied_seats"`
	FreeSeats         int32     `json:"free_seats"`
	RegisteredReaders int64     `json:"registered_readers"`
	ActiveReaders     int64     `json:"active_readers"`
}

// listReadingHalls возвращает список всех читальных залов
func (h *Handler) listReadingHalls(c *fiber.Ctx) error {
	halls, err := h.repo.ListReadingHalls(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to list reading halls")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reading halls.", err.Error())
	}

	response := make([]ReadingHallResponse, len(halls))
	for i, hall := range halls {
		freeSeats := 0
		if hall.OccupiedSeats != nil {
			freeSeats = hall.TotalSeats - *hall.OccupiedSeats
		} else {
			freeSeats = hall.TotalSeats
		}

		response[i] = ReadingHallResponse{
			ID:             hall.ID,
			LibraryName:    hall.LibraryName,
			HallName:       hall.HallName,
			Specialization: hall.Specialization,
			TotalSeats:     hall.TotalSeats,
			OccupiedSeats:  hall.OccupiedSeats,
			FreeSeats:      &freeSeats,
			CreatedAt:      hall.CreatedAt,
		}
	}

	return c.JSON(fiber.Map{"halls": response})
}

// getHallStatistics возвращает статистику по всем залам
func (h *Handler) getHallStatistics(c *fiber.Ctx) error {
	statistics, err := h.repo.GetHallStatistics(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get hall statistics")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve hall statistics.", err.Error())
	}

	response := make([]HallStatisticsResponse, len(statistics))
	for i, stat := range statistics {
		response[i] = HallStatisticsResponse{
			ID:                stat.ID,
			LibraryName:       stat.LibraryName,
			HallName:          stat.HallName,
			Specialization:    stat.Specialization,
			TotalSeats:        stat.TotalSeats,
			OccupiedSeats:     stat.OccupiedSeats,
			FreeSeats:         stat.FreeSeats,
			RegisteredReaders: stat.RegisteredReaders,
			ActiveReaders:     stat.ActiveReaders,
		}
	}

	return c.JSON(fiber.Map{"statistics": response})
}

// getReadingHallByID возвращает читальный зал по ID
func (h *Handler) getReadingHallByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	hallID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format.")
	}

	hall, err := h.repo.GetReadingHallByID(c.Context(), hallID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reading hall not found.")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to get reading hall by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reading hall.", err.Error())
	}

	freeSeats := 0
	if hall.OccupiedSeats != nil {
		freeSeats = hall.TotalSeats - *hall.OccupiedSeats
	} else {
		freeSeats = hall.TotalSeats
	}

	response := ReadingHallResponse{
		ID:             hall.ID,
		LibraryName:    hall.LibraryName,
		HallName:       hall.HallName,
		Specialization: hall.Specialization,
		TotalSeats:     hall.TotalSeats,
		OccupiedSeats:  hall.OccupiedSeats,
		FreeSeats:      &freeSeats,
		CreatedAt:      hall.CreatedAt,
	}

	return c.JSON(response)
}

// createReadingHall создает новый читальный зал
func (h *Handler) createReadingHall(c *fiber.Ctx) error {
	var req CreateReadingHallRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация обязательных полей
	if req.LibraryName == "" {
		return httperr.New(fiber.StatusBadRequest, "Library name is required.")
	}
	if req.HallName == "" {
		return httperr.New(fiber.StatusBadRequest, "Hall name is required.")
	}
	if req.TotalSeats <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Total seats must be greater than 0.")
	}

	// Создание читального зала
	hall, err := h.repo.CreateReadingHall(c.Context(), postgres.CreateReadingHallParams{
		LibraryName:    req.LibraryName,
		HallName:       req.HallName,
		Specialization: req.Specialization,
		TotalSeats:     req.TotalSeats,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return httperr.New(fiber.StatusConflict, "Reading hall with this name already exists in the library.")
		}
		log.Error().Err(err).Msg("Failed to create reading hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create reading hall.", err.Error())
	}

	freeSeats := hall.TotalSeats
	if hall.OccupiedSeats != nil {
		freeSeats = hall.TotalSeats - *hall.OccupiedSeats
	}

	response := ReadingHallResponse{
		ID:             hall.ID,
		LibraryName:    hall.LibraryName,
		HallName:       hall.HallName,
		Specialization: hall.Specialization,
		TotalSeats:     hall.TotalSeats,
		OccupiedSeats:  hall.OccupiedSeats,
		FreeSeats:      &freeSeats,
		CreatedAt:      hall.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateReadingHall обновляет информацию о читальном зале
func (h *Handler) updateReadingHall(c *fiber.Ctx) error {
	idStr := c.Params("id")
	hallID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format.")
	}

	var req UpdateReadingHallRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию о зале для сохранения неизменяемых полей
	existingHall, err := h.repo.GetReadingHallByID(c.Context(), hallID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reading hall not found.")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to get reading hall for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reading hall.", err.Error())
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateReadingHallParams{
		HallID: hallID,
	}

	// Устанавливаем значения для обновления или сохраняем существующие
	if req.LibraryName != nil {
		updateParams.LibraryName = *req.LibraryName
	} else {
		updateParams.LibraryName = existingHall.LibraryName
	}

	if req.HallName != nil {
		updateParams.HallName = *req.HallName
	} else {
		updateParams.HallName = existingHall.HallName
	}

	if req.Specialization != nil {
		updateParams.Specialization = req.Specialization
	} else {
		updateParams.Specialization = existingHall.Specialization
	}

	if req.TotalSeats != nil {
		if *req.TotalSeats <= 0 {
			return httperr.New(fiber.StatusBadRequest, "Total seats must be greater than 0.")
		}
		updateParams.TotalSeats = *req.TotalSeats
	} else {
		updateParams.TotalSeats = existingHall.TotalSeats
	}

	if req.OccupiedSeats != nil {
		if *req.OccupiedSeats < 0 {
			return httperr.New(fiber.StatusBadRequest, "Occupied seats cannot be negative.")
		}
		if *req.OccupiedSeats > updateParams.TotalSeats {
			return httperr.New(fiber.StatusBadRequest, "Occupied seats cannot exceed total seats.")
		}
		updateParams.OccupiedSeats = req.OccupiedSeats
	} else {
		updateParams.OccupiedSeats = existingHall.OccupiedSeats
	}

	// Обновление читального зала
	updatedHall, err := h.repo.UpdateReadingHall(c.Context(), updateParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return httperr.New(fiber.StatusConflict, "Reading hall with this name already exists in the library.")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to update reading hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update reading hall.", err.Error())
	}

	freeSeats := updatedHall.TotalSeats
	if updatedHall.OccupiedSeats != nil {
		freeSeats = updatedHall.TotalSeats - *updatedHall.OccupiedSeats
	}

	response := ReadingHallResponse{
		ID:             updatedHall.ID,
		LibraryName:    updatedHall.LibraryName,
		HallName:       updatedHall.HallName,
		Specialization: updatedHall.Specialization,
		TotalSeats:     updatedHall.TotalSeats,
		OccupiedSeats:  updatedHall.OccupiedSeats,
		FreeSeats:      &freeSeats,
		CreatedAt:      updatedHall.CreatedAt,
	}

	return c.JSON(response)
}

// updateHallOccupancy обновляет информацию о занятости зала
func (h *Handler) updateHallOccupancy(c *fiber.Ctx) error {
	idStr := c.Params("id")
	hallID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format.")
	}

	var req UpdateHallOccupancyRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Проверяем, что зал существует и получаем информацию о нем для валидации
	existingHall, err := h.repo.GetReadingHallByID(c.Context(), hallID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reading hall not found.")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to get reading hall for occupancy update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reading hall.", err.Error())
	}

	// Валидация занятости
	if req.OccupiedSeats != nil {
		if *req.OccupiedSeats < 0 {
			return httperr.New(fiber.StatusBadRequest, "Occupied seats cannot be negative.")
		}
		if *req.OccupiedSeats > existingHall.TotalSeats {
			return httperr.New(fiber.StatusBadRequest, "Occupied seats cannot exceed total seats.")
		}
	}

	// Обновление занятости зала
	err = h.repo.UpdateHallOccupancy(c.Context(), postgres.UpdateHallOccupancyParams{
		HallID:        hallID,
		OccupiedSeats: req.OccupiedSeats,
	})
	if err != nil {
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to update hall occupancy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update hall occupancy.", err.Error())
	}

	// Получаем обновленную информацию
	updatedHall, err := h.repo.GetReadingHallByID(c.Context(), hallID)
	if err != nil {
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to get updated hall info")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve updated hall info.", err.Error())
	}

	freeSeats := updatedHall.TotalSeats
	if updatedHall.OccupiedSeats != nil {
		freeSeats = updatedHall.TotalSeats - *updatedHall.OccupiedSeats
	}

	response := ReadingHallResponse{
		ID:             updatedHall.ID,
		LibraryName:    updatedHall.LibraryName,
		HallName:       updatedHall.HallName,
		Specialization: updatedHall.Specialization,
		TotalSeats:     updatedHall.TotalSeats,
		OccupiedSeats:  updatedHall.OccupiedSeats,
		FreeSeats:      &freeSeats,
		CreatedAt:      updatedHall.CreatedAt,
	}

	return c.JSON(response)
}

// deleteReadingHall удаляет читальный зал
func (h *Handler) deleteReadingHall(c *fiber.Ctx) error {
	idStr := c.Params("id")
	hallID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format.")
	}

	// Проверяем, что зал существует
	_, err = h.repo.GetReadingHallByID(c.Context(), hallID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reading hall not found.")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to get reading hall for deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reading hall.", err.Error())
	}

	// Проверяем, есть ли активные читатели в этом зале
	readersCount, err := h.repo.CountReadersByHall(c.Context(), &hallID)
	if err != nil {
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to count readers in hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to check hall usage.", err.Error())
	}

	if readersCount > 0 {
		return httperr.New(fiber.StatusConflict, "Cannot delete reading hall with active readers. Please reassign readers first.")
	}

	// Удаление читального зала
	err = h.repo.DeleteReadingHall(c.Context(), hallID)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusConflict, "Cannot delete reading hall that is referenced by other records.")
		}
		log.Error().Err(err).Str("hallID", idStr).Msg("Failed to delete reading hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to delete reading hall.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

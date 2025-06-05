package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

// Request/Response structs для readers
type CreateReaderRequest struct {
	UserID        *uuid.UUID `json:"user_id"`
	TicketNumber  string     `json:"ticket_number" validate:"required,min=3,max=20"`
	FullName      string     `json:"full_name" validate:"required,min=2,max=100"`
	BirthDate     time.Time  `json:"birth_date" validate:"required"`
	Phone         *string    `json:"phone" validate:"omitempty,min=10,max=20"`
	Education     *string    `json:"education" validate:"omitempty,max=100"`
	ReadingHallID *uuid.UUID `json:"reading_hall_id"`
}

type UpdateReaderRequest struct {
	FullName      *string    `json:"full_name" validate:"omitempty,min=2,max=100"`
	Phone         *string    `json:"phone" validate:"omitempty,min=10,max=20"`
	Education     *string    `json:"education" validate:"omitempty,max=100"`
	ReadingHallID *uuid.UUID `json:"reading_hall_id"`
}

type ReaderResponse struct {
	ID               uuid.UUID  `json:"id"`
	UserID           *uuid.UUID `json:"user_id"`
	TicketNumber     string     `json:"ticket_number"`
	FullName         string     `json:"full_name"`
	BirthDate        time.Time  `json:"birth_date"`
	Phone            *string    `json:"phone"`
	Education        *string    `json:"education"`
	ReadingHallID    *uuid.UUID `json:"reading_hall_id"`
	RegistrationDate *time.Time `json:"registration_date"`
	IsActive         *bool      `json:"is_active"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	Username         *string    `json:"username,omitempty"`
	Email            *string    `json:"email,omitempty"`
	HallName         *string    `json:"hall_name,omitempty"`
	Specialization   *string    `json:"specialization,omitempty"`
}

type ReadersListResponse struct {
	Readers []ReaderResponse `json:"readers"`
	Total   int64            `json:"total"`
	Limit   int32            `json:"limit"`
	Offset  int32            `json:"offset"`
}

// listReaders возвращает список всех читателей с пагинацией
func (h *Handler) listReaders(c *fiber.Ctx) error {
	// Парсинг параметров пагинации
	limit := int32(20)
	offset := int32(0)

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = int32(parsedLimit)
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = int32(parsedOffset)
		}
	}

	// Получение читателей из БД
	readers, err := h.repo.ListAllReaders(c.Context(), postgres.ListAllReadersParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to list readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve readers.", err.Error())
	}

	// Получение общего количества читателей
	total, err := h.repo.CountReaders(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to count readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to count readers.", err.Error())
	}

	// Преобразование в response структуры
	response := ReadersListResponse{
		Readers: make([]ReaderResponse, len(readers)),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}

	for i, reader := range readers {
		response.Readers[i] = ReaderResponse{
			ID:               reader.ID,
			UserID:           reader.UserID,
			TicketNumber:     reader.TicketNumber,
			FullName:         reader.FullName,
			BirthDate:        reader.BirthDate,
			Phone:            reader.Phone,
			Education:        reader.Education,
			ReadingHallID:    reader.ReadingHallID,
			RegistrationDate: reader.RegistrationDate,
			IsActive:         reader.IsActive,
			CreatedAt:        reader.CreatedAt,
			UpdatedAt:        reader.UpdatedAt,
			Username:         &reader.Username,
			Email:            &reader.Email,
			HallName:         reader.HallName,
		}
	}

	return c.JSON(response)
}

// getCurrentReader возвращает информацию о текущем читателе
func (h *Handler) getCurrentReader(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID.")
	}

	reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader profile not found.")
		}
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get current reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	response := ReaderResponse{
		ID:               reader.ID,
		UserID:           reader.UserID,
		TicketNumber:     reader.TicketNumber,
		FullName:         reader.FullName,
		BirthDate:        reader.BirthDate,
		Phone:            reader.Phone,
		Education:        reader.Education,
		ReadingHallID:    reader.ReadingHallID,
		RegistrationDate: reader.RegistrationDate,
		IsActive:         reader.IsActive,
		CreatedAt:        reader.CreatedAt,
		UpdatedAt:        reader.UpdatedAt,
		HallName:         reader.HallName,
		Specialization:   reader.Specialization,
	}

	return c.JSON(response)
}

// getReaderByID возвращает читателя по ID
func (h *Handler) getReaderByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	readerID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format.")
	}

	reader, err := h.repo.GetReaderByID(c.Context(), readerID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to get reader by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader.", err.Error())
	}

	response := ReaderResponse{
		ID:               reader.ID,
		UserID:           reader.UserID,
		TicketNumber:     reader.TicketNumber,
		FullName:         reader.FullName,
		BirthDate:        reader.BirthDate,
		Phone:            reader.Phone,
		Education:        reader.Education,
		ReadingHallID:    reader.ReadingHallID,
		RegistrationDate: reader.RegistrationDate,
		IsActive:         reader.IsActive,
		CreatedAt:        reader.CreatedAt,
		UpdatedAt:        reader.UpdatedAt,
		Username:         &reader.Username,
		Email:            &reader.Email,
		HallName:         reader.HallName,
		Specialization:   reader.Specialization,
	}

	return c.JSON(response)
}

// getReaderByTicket возвращает читателя по номеру билета
func (h *Handler) getReaderByTicket(c *fiber.Ctx) error {
	ticketNumber := c.Params("ticket")
	if ticketNumber == "" {
		return httperr.New(fiber.StatusBadRequest, "Ticket number is required.")
	}

	reader, err := h.repo.GetReaderByTicketNumber(c.Context(), ticketNumber)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader with this ticket number not found.")
		}
		log.Error().Err(err).Str("ticketNumber", ticketNumber).Msg("Failed to get reader by ticket")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader.", err.Error())
	}

	response := ReaderResponse{
		ID:               reader.ID,
		UserID:           reader.UserID,
		TicketNumber:     reader.TicketNumber,
		FullName:         reader.FullName,
		BirthDate:        reader.BirthDate,
		Phone:            reader.Phone,
		Education:        reader.Education,
		ReadingHallID:    reader.ReadingHallID,
		RegistrationDate: reader.RegistrationDate,
		IsActive:         reader.IsActive,
		CreatedAt:        reader.CreatedAt,
		UpdatedAt:        reader.UpdatedAt,
		Username:         &reader.Username,
		Email:            &reader.Email,
		HallName:         reader.HallName,
		Specialization:   reader.Specialization,
	}

	return c.JSON(response)
}

// getReadersByHall возвращает читателей по залу
func (h *Handler) getReadersByHall(c *fiber.Ctx) error {
	hallIDStr := c.Params("hallId")
	hallID, err := uuid.Parse(hallIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall ID format.")
	}

	readers, err := h.repo.ListReadersByHall(c.Context(), &hallID)
	if err != nil {
		log.Error().Err(err).Str("hallID", hallIDStr).Msg("Failed to get readers by hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve readers.", err.Error())
	}

	response := make([]ReaderResponse, len(readers))
	for i, reader := range readers {
		response[i] = ReaderResponse{
			ID:               reader.ID,
			UserID:           reader.UserID,
			TicketNumber:     reader.TicketNumber,
			FullName:         reader.FullName,
			BirthDate:        reader.BirthDate,
			Phone:            reader.Phone,
			Education:        reader.Education,
			ReadingHallID:    reader.ReadingHallID,
			RegistrationDate: reader.RegistrationDate,
			IsActive:         reader.IsActive,
			CreatedAt:        reader.CreatedAt,
			UpdatedAt:        reader.UpdatedAt,
			Username:         &reader.Username,
			Email:            &reader.Email,
		}
	}

	return c.JSON(fiber.Map{"readers": response})
}

// searchReaders ищет читателей по имени
func (h *Handler) searchReaders(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return httperr.New(fiber.StatusBadRequest, "Search query parameter 'q' is required.")
	}

	readers, err := h.repo.SearchReadersByName(c.Context(), &query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to search readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search readers.", err.Error())
	}

	response := make([]ReaderResponse, len(readers))
	for i, reader := range readers {
		response[i] = ReaderResponse{
			ID:               reader.ID,
			UserID:           reader.UserID,
			TicketNumber:     reader.TicketNumber,
			FullName:         reader.FullName,
			BirthDate:        reader.BirthDate,
			Phone:            reader.Phone,
			Education:        reader.Education,
			ReadingHallID:    reader.ReadingHallID,
			RegistrationDate: reader.RegistrationDate,
			IsActive:         reader.IsActive,
			CreatedAt:        reader.CreatedAt,
			UpdatedAt:        reader.UpdatedAt,
			Username:         &reader.Username,
			Email:            &reader.Email,
			HallName:         reader.HallName,
		}
	}

	return c.JSON(fiber.Map{"readers": response})
}

// createReader создает нового читателя
func (h *Handler) createReader(c *fiber.Ctx) error {
	var req CreateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация обязательных полей
	if req.TicketNumber == "" {
		return httperr.New(fiber.StatusBadRequest, "Ticket number is required.")
	}
	if req.FullName == "" {
		return httperr.New(fiber.StatusBadRequest, "Full name is required.")
	}
	if req.BirthDate.IsZero() {
		return httperr.New(fiber.StatusBadRequest, "Birth date is required.")
	}

	// Создание читателя
	reader, err := h.repo.CreateReader(c.Context(), postgres.CreateReaderParams{
		UserID:        req.UserID,
		TicketNumber:  req.TicketNumber,
		FullName:      req.FullName,
		BirthDate:     req.BirthDate,
		Phone:         req.Phone,
		Education:     req.Education,
		ReadingHallID: req.ReadingHallID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return httperr.New(fiber.StatusConflict, "Reader with this ticket number already exists.")
		}
		log.Error().Err(err).Msg("Failed to create reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create reader.", err.Error())
	}

	response := ReaderResponse{
		ID:               reader.ID,
		UserID:           reader.UserID,
		TicketNumber:     reader.TicketNumber,
		FullName:         reader.FullName,
		BirthDate:        reader.BirthDate,
		Phone:            reader.Phone,
		Education:        reader.Education,
		ReadingHallID:    reader.ReadingHallID,
		RegistrationDate: reader.RegistrationDate,
		IsActive:         reader.IsActive,
		CreatedAt:        reader.CreatedAt,
		UpdatedAt:        reader.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateReader обновляет информацию о читателе
func (h *Handler) updateReader(c *fiber.Ctx) error {
	idStr := c.Params("id")
	readerID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format.")
	}

	var req UpdateReaderRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию о читателе для сохранения неизменяемых полей
	existingReader, err := h.repo.GetReaderByID(c.Context(), readerID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to get reader for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader.", err.Error())
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateReaderParams{
		ReaderID: readerID,
	}

	// Устанавливаем значения для обновления или сохраняем существующие
	if req.FullName != nil {
		updateParams.FullName = *req.FullName
	} else {
		updateParams.FullName = existingReader.FullName
	}

	if req.Phone != nil {
		updateParams.Phone = req.Phone
	} else {
		updateParams.Phone = existingReader.Phone
	}

	if req.Education != nil {
		updateParams.Education = req.Education
	} else {
		updateParams.Education = existingReader.Education
	}

	if req.ReadingHallID != nil {
		updateParams.ReadingHallID = req.ReadingHallID
	} else {
		updateParams.ReadingHallID = existingReader.ReadingHallID
	}

	// Обновление читателя
	updatedReader, err := h.repo.UpdateReader(c.Context(), updateParams)
	if err != nil {
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to update reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update reader.", err.Error())
	}

	response := ReaderResponse{
		ID:               updatedReader.ID,
		UserID:           updatedReader.UserID,
		TicketNumber:     updatedReader.TicketNumber,
		FullName:         updatedReader.FullName,
		BirthDate:        updatedReader.BirthDate,
		Phone:            updatedReader.Phone,
		Education:        updatedReader.Education,
		ReadingHallID:    updatedReader.ReadingHallID,
		RegistrationDate: updatedReader.RegistrationDate,
		IsActive:         updatedReader.IsActive,
		CreatedAt:        updatedReader.CreatedAt,
		UpdatedAt:        updatedReader.UpdatedAt,
	}

	return c.JSON(response)
}

// deactivateReader деактивирует читателя
func (h *Handler) deactivateReader(c *fiber.Ctx) error {
	idStr := c.Params("id")
	readerID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format.")
	}

	// Проверяем, что читатель существует
	_, err = h.repo.GetReaderByID(c.Context(), readerID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to get reader for deactivation")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader.", err.Error())
	}

	// Деактивация читателя
	err = h.repo.DeactivateReader(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", idStr).Msg("Failed to deactivate reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to deactivate reader.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

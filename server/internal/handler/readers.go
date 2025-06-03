package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type CreateReaderRequest struct {
	FullName     string    `json:"full_name"`
	TicketNumber string    `json:"ticket_number"`
	BirthDate    time.Time `json:"birth_date"`
	Phone        *string   `json:"phone"`
	Email        *string   `json:"email"`
	Education    *string   `json:"education"`
	HallID       int       `json:"hall_id"`
}

type UpdateReaderRequest struct {
	FullName  string  `json:"full_name"`
	Phone     *string `json:"phone"`
	Email     *string `json:"email"`
	Education *string `json:"education"`
	HallID    int     `json:"hall_id"`
}

type UpdateReaderStatusRequest struct {
	Status string `json:"status"`
}

type SearchReadersByNameRequest struct {
	SearchName string `json:"search_name"`
	PageOffset int32  `json:"page_offset"`
	PageLimit  int32  `json:"page_limit"`
}

type GetAllReadersRequest struct {
	PageOffset int32 `json:"page_offset"`
	PageLimit  int32 `json:"page_limit"`
}

// createReader создает нового читателя
func (h *Handler) createReader(c *fiber.Ctx) error {
	req := new(CreateReaderRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: full_name, ticket_number, birth_date, hall_id
	// TODO: Validate full_name min length 2, max length 100
	// TODO: Validate ticket_number is unique and matches format
	// TODO: Validate birth_date is not in future and reader is at least 6 years old
	// TODO: Validate phone format if provided
	// TODO: Validate email format if provided
	// TODO: Validate hall_id exists

	params := postgres.CreateReaderParams{
		FullName:     req.FullName,
		TicketNumber: req.TicketNumber,
		BirthDate:    req.BirthDate,
		Phone:        req.Phone,
		Email:        req.Email,
		Education:    req.Education,
		HallID:       req.HallID,
	}

	reader, err := h.repo.CreateReader(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create reader")
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Reader with this ticket number already exists.")
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid hall ID.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create reader.")
	}

	return c.Status(fiber.StatusCreated).JSON(reader)
}

// getReaderByID получает читателя по ID
func (h *Handler) getReaderByID(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	reader, err := h.repo.GetReaderByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Int64("readerID", id).Msg("Failed to get reader by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader.")
	}

	return c.JSON(reader)
}

// getReaderByTicket получает читателя по номеру билета
func (h *Handler) getReaderByTicket(c *fiber.Ctx) error {
	ticketNumber := c.Params("ticket")
	if ticketNumber == "" {
		return httperr.New(fiber.StatusBadRequest, "Ticket number is required.")
	}

	// TODO: Validate ticket_number format

	reader, err := h.repo.GetReaderByTicket(c.Context(), ticketNumber)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Str("ticketNumber", ticketNumber).Msg("Failed to get reader by ticket")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader.")
	}

	return c.JSON(reader)
}

// getAllReaders получает список всех читателей с пагинацией
func (h *Handler) getAllReaders(c *fiber.Ctx) error {
	req := new(GetAllReadersRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0

	if req.PageLimit == 0 {
		req.PageLimit = 20 // default limit
	}

	params := postgres.GetAllReadersParams{
		PageOffset: req.PageOffset,
		PageLimit:  req.PageLimit,
	}

	readers, err := h.repo.GetAllReaders(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve readers.")
	}

	if readers == nil {
		readers = []*postgres.GetAllReadersRow{}
	}

	return c.JSON(readers)
}

// getActiveReaders получает список активных читателей
func (h *Handler) getActiveReaders(c *fiber.Ctx) error {
	limit := int32(50) // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// TODO: Validate limit > 0 and <= 100

	readers, err := h.repo.GetActiveReaders(c.Context(), limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get active readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve active readers.")
	}

	if readers == nil {
		readers = []*postgres.GetActiveReadersRow{}
	}

	return c.JSON(readers)
}

// searchReadersByName ищет читателей по имени
func (h *Handler) searchReadersByName(c *fiber.Ctx) error {
	req := new(SearchReadersByNameRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate search_name is not empty and has min length 2
	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0

	if req.PageLimit == 0 {
		req.PageLimit = 20 // default limit
	}

	params := postgres.SearchReadersByNameParams{
		SearchName: &req.SearchName,
		PageOffset: req.PageOffset,
		PageLimit:  req.PageLimit,
	}

	readers, err := h.repo.SearchReadersByName(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Str("searchName", req.SearchName).Msg("Failed to search readers by name")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search readers.")
	}

	if readers == nil {
		readers = []*postgres.SearchReadersByNameRow{}
	}

	return c.JSON(readers)
}

// updateReader обновляет информацию о читателе
func (h *Handler) updateReader(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateReaderRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: full_name, hall_id
	// TODO: Validate full_name min length 2, max length 100
	// TODO: Validate phone format if provided
	// TODO: Validate email format if provided
	// TODO: Validate hall_id exists

	params := postgres.UpdateReaderParams{
		ReaderID:  id,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Email:     req.Email,
		Education: req.Education,
		HallID:    req.HallID,
	}

	reader, err := h.repo.UpdateReader(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid hall ID.")
		}
		log.Error().Err(err).Int64("readerID", id).Msg("Failed to update reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update reader.")
	}

	return c.JSON(reader)
}

// updateReaderStatus обновляет статус читателя
func (h *Handler) updateReaderStatus(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateReaderStatusRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate status is one of: active, suspended, blocked, inactive

	params := postgres.UpdateReaderStatusParams{
		ReaderID: id,
		Status:   req.Status,
	}

	err = h.repo.UpdateReaderStatus(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Int64("readerID", id).Msg("Failed to update reader status")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update reader status.")
	}

	return c.JSON(fiber.Map{"message": "Reader status updated successfully"})
}

// updateReaderDebt обновляет задолженность читателя
func (h *Handler) updateReaderDebt(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.UpdateReaderDebt(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Int64("readerID", id).Msg("Failed to update reader debt")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update reader debt.")
	}

	return c.JSON(fiber.Map{"message": "Reader debt updated successfully"})
}

// getReaderStatistics получает статистику читателя
func (h *Handler) getReaderStatistics(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	statistics, err := h.repo.GetReaderStatistics(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Int64("readerID", id).Msg("Failed to get reader statistics")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader statistics.")
	}

	return c.JSON(statistics)
}

// getReaderFavoriteCategories получает любимые категории читателя
func (h *Handler) getReaderFavoriteCategories(c *fiber.Ctx) error {
	readerIDStr := c.Params("id")
	if readerIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Reader ID is required.")
	}

	readerID, err := strconv.Atoi(readerIDStr)
	if err != nil || readerID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID.")
	}

	// TODO: Validate reader_id exists

	categories, err := h.repo.GetReaderFavoriteCategories(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Int("readerID", readerID).Msg("Failed to get reader favorite categories")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve favorite categories.")
	}

	if categories == nil {
		categories = []*postgres.GetReaderFavoriteCategoriesRow{}
	}

	return c.JSON(categories)
}

// getReadersCount получает общее количество читателей
func (h *Handler) getReadersCount(c *fiber.Ctx) error {
	count, err := h.repo.GetReadersCount(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get readers count")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve readers count.")
	}

	return c.JSON(fiber.Map{"total_readers": count})
}

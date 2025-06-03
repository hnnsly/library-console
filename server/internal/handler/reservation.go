package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type CreateReservationRequest struct {
	BookID   int `json:"book_id"`
	ReaderID int `json:"reader_id"`
}

// createReservation создает новое бронирование
func (h *Handler) createReservation(c *fiber.Ctx) error {
	req := new(CreateReservationRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: book_id, reader_id
	// TODO: Validate book_id and reader_id exist
	// TODO: Validate book is not currently available
	// TODO: Validate reader doesn't already have reservation for this book
	// TODO: Validate reader doesn't exceed max reservations allowed
	// TODO: Validate reader status is active

	params := postgres.CreateReservationParams{
		BookID:   req.BookID,
		ReaderID: req.ReaderID,
	}

	reservation, err := h.repo.CreateReservation(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create reservation")
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid book or reader ID.")
		}
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Reader already has a reservation for this book.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create reservation.")
	}

	return c.Status(fiber.StatusCreated).JSON(reservation)
}

// getReaderReservations получает бронирования читателя
func (h *Handler) getReaderReservations(c *fiber.Ctx) error {
	readerIDStr := c.Params("reader_id")
	if readerIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Reader ID is required.")
	}

	readerID, err := strconv.Atoi(readerIDStr)
	if err != nil || readerID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID.")
	}

	// TODO: Validate reader_id exists

	reservations, err := h.repo.GetReaderReservations(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Int("readerID", readerID).Msg("Failed to get reader reservations")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reservations.")
	}

	if reservations == nil {
		reservations = []*postgres.GetReaderReservationsRow{}
	}

	return c.JSON(reservations)
}

// getBookQueue получает очередь бронирований для книги
func (h *Handler) getBookQueue(c *fiber.Ctx) error {
	bookIDStr := c.Params("book_id")
	if bookIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Book ID is required.")
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID.")
	}

	// TODO: Validate book_id exists

	queue, err := h.repo.GetBookQueue(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Int("bookID", bookID).Msg("Failed to get book queue")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book queue.")
	}

	if queue == nil {
		queue = []*postgres.GetBookQueueRow{}
	}

	return c.JSON(queue)
}

// getExpiredReservations получает просроченные бронирования
func (h *Handler) getExpiredReservations(c *fiber.Ctx) error {
	reservations, err := h.repo.GetExpiredReservations(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get expired reservations")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve expired reservations.")
	}

	if reservations == nil {
		reservations = []*postgres.GetExpiredReservationsRow{}
	}

	return c.JSON(reservations)
}

// fulfillReservation исполняет бронирование (когда книга выдается)
func (h *Handler) fulfillReservation(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.FulfillReservation(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Reservation not found.")
		}
		log.Error().Err(err).Int64("reservationID", id).Msg("Failed to fulfill reservation")
		return httperr.New(fiber.StatusInternalServerError, "Failed to fulfill reservation.")
	}

	return c.JSON(fiber.Map{"message": "Reservation successfully fulfilled"})
}

// cancelReservation отменяет бронирование
func (h *Handler) cancelReservation(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.CancelReservation(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Reservation not found.")
		}
		log.Error().Err(err).Int64("reservationID", id).Msg("Failed to cancel reservation")
		return httperr.New(fiber.StatusInternalServerError, "Failed to cancel reservation.")
	}

	return c.JSON(fiber.Map{"message": "Reservation successfully cancelled"})
}

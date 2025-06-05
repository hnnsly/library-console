package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

// Request/Response structs для ratings
type CreateRatingRequest struct {
	BookID     uuid.UUID  `json:"book_id" validate:"required"`
	Rating     int        `json:"rating" validate:"required,min=1,max=5"`
	Review     *string    `json:"review" validate:"omitempty,max=2000"`
	RatingDate *time.Time `json:"rating_date"`
}

type UpdateRatingRequest struct {
	Rating *int    `json:"rating" validate:"omitempty,min=1,max=5"`
	Review *string `json:"review" validate:"omitempty,max=2000"`
}

type RatingResponse struct {
	ID         uuid.UUID  `json:"id"`
	BookID     uuid.UUID  `json:"book_id"`
	ReaderID   uuid.UUID  `json:"reader_id"`
	Rating     int        `json:"rating"`
	Review     *string    `json:"review"`
	RatingDate *time.Time `json:"rating_date"`
	CreatedAt  *time.Time `json:"created_at"`
	// Extended fields
	ReaderName *string `json:"reader_name,omitempty"`
	BookTitle  *string `json:"book_title,omitempty"`
}

type BookAverageRatingResponse struct {
	BookID       uuid.UUID       `json:"book_id"`
	AvgRating    decimal.Decimal `json:"avg_rating"`
	TotalRatings int64           `json:"total_ratings"`
}

type TopRatedBookResponse struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Authors     string          `json:"authors"`
	AvgRating   decimal.Decimal `json:"avg_rating"`
	RatingCount int64           `json:"rating_count"`
}

// getMyRatings возвращает рейтинги текущего читателя
func (h *Handler) getMyRatings(c *fiber.Ctx) error {
	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID.")
	}

	// Получаем информацию о читателе
	reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader profile not found.")
		}
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get reader for ratings")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	return h.getReaderRatingsInternal(c, reader.ID)
}

// getReaderRatings возвращает рейтинги читателя по ID
func (h *Handler) getReaderRatings(c *fiber.Ctx) error {
	readerIDStr := c.Params("readerId")
	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format.")
	}

	// Проверка доступа: читатель может видеть только свои рейтинги
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != readerID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own ratings.")
		}
	}

	return h.getReaderRatingsInternal(c, readerID)
}

// getReaderRatingsInternal - внутренний метод для получения рейтингов читателя
func (h *Handler) getReaderRatingsInternal(c *fiber.Ctx, readerID uuid.UUID) error {
	ratings, err := h.repo.GetReaderRatings(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", readerID.String()).Msg("Failed to get reader ratings")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve ratings.", err.Error())
	}

	response := make([]RatingResponse, len(ratings))
	for i, rating := range ratings {
		response[i] = RatingResponse{
			ID:         rating.ID,
			BookID:     rating.BookID,
			ReaderID:   rating.ReaderID,
			Rating:     rating.Rating,
			Review:     rating.Review,
			RatingDate: rating.RatingDate,
			CreatedAt:  rating.CreatedAt,
			BookTitle:  &rating.BookTitle,
		}
	}

	return c.JSON(fiber.Map{"ratings": response})
}

// getBookRatings возвращает все рейтинги книги
func (h *Handler) getBookRatings(c *fiber.Ctx) error {
	bookIDStr := c.Params("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	ratings, err := h.repo.GetBookRatings(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to get book ratings")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book ratings.", err.Error())
	}

	response := make([]RatingResponse, len(ratings))
	for i, rating := range ratings {
		response[i] = RatingResponse{
			ID:         rating.ID,
			BookID:     rating.BookID,
			ReaderID:   rating.ReaderID,
			Rating:     rating.Rating,
			Review:     rating.Review,
			RatingDate: rating.RatingDate,
			CreatedAt:  rating.CreatedAt,
			ReaderName: &rating.ReaderName,
		}
	}

	return c.JSON(fiber.Map{"ratings": response})
}

// getBookAverageRating возвращает средний рейтинг книги
func (h *Handler) getBookAverageRating(c *fiber.Ctx) error {
	bookIDStr := c.Params("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	avgRating, err := h.repo.GetBookAverageRating(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to get book average rating")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book average rating.", err.Error())
	}

	// Конвертируем interface{} в decimal.Decimal
	rating, _ := decimal.Parse(avgRating.AvgRating.(string))

	response := BookAverageRatingResponse{
		BookID:       bookID,
		AvgRating:    rating,
		TotalRatings: avgRating.TotalRatings,
	}

	return c.JSON(response)
}

// getMyBookRating возвращает рейтинг текущего читателя для книги
func (h *Handler) getMyBookRating(c *fiber.Ctx) error {
	bookIDStr := c.Params("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID.")
	}

	// Получаем информацию о читателе
	reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader profile not found.")
		}
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get reader for rating")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	rating, err := h.repo.GetReaderBookRating(c.Context(), postgres.GetReaderBookRatingParams{
		BookID:   bookID,
		ReaderID: reader.ID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Rating not found.")
		}
		log.Error().Err(err).Str("bookID", bookIDStr).Str("readerID", reader.ID.String()).Msg("Failed to get reader book rating")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve rating.", err.Error())
	}

	response := RatingResponse{
		ID:         rating.ID,
		BookID:     rating.BookID,
		ReaderID:   rating.ReaderID,
		Rating:     rating.Rating,
		Review:     rating.Review,
		RatingDate: rating.RatingDate,
		CreatedAt:  rating.CreatedAt,
	}

	return c.JSON(response)
}

// getRatingByID возвращает рейтинг по ID
func (h *Handler) getRatingByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	ratingID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid rating ID format.")
	}

	rating, err := h.repo.GetBookRatingByID(c.Context(), ratingID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Rating not found.")
		}
		log.Error().Err(err).Str("ratingID", idStr).Msg("Failed to get rating by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve rating.", err.Error())
	}

	// Проверка доступа для читателей
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != rating.ReaderID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own ratings.")
		}
	}

	response := RatingResponse{
		ID:         rating.ID,
		BookID:     rating.BookID,
		ReaderID:   rating.ReaderID,
		Rating:     rating.Rating,
		Review:     rating.Review,
		RatingDate: rating.RatingDate,
		CreatedAt:  rating.CreatedAt,
		ReaderName: &rating.ReaderName,
		BookTitle:  &rating.BookTitle,
	}

	return c.JSON(response)
}

// createRating создает новый рейтинг книги
func (h *Handler) createRating(c *fiber.Ctx) error {
	var req CreateRatingRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация
	if req.Rating < 1 || req.Rating > 5 {
		return httperr.New(fiber.StatusBadRequest, "Rating must be between 1 and 5.")
	}

	userIDStr := c.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid user ID.")
	}

	// Получаем информацию о читателе
	reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader profile not found.")
		}
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get reader for rating creation")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	// Проверяем, что книга существует
	_, err = h.repo.GetBookByID(c.Context(), req.BookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", req.BookID.String()).Msg("Failed to verify book existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify book.", err.Error())
	}

	// Если дата рейтинга не указана, используем текущую дату
	if req.RatingDate == nil {
		now := time.Now()
		req.RatingDate = &now
	}

	// Создание или обновление рейтинга (используем UPSERT)
	rating, err := h.repo.CreateBookRating(c.Context(), postgres.CreateBookRatingParams{
		BookID:     req.BookID,
		ReaderID:   reader.ID,
		Rating:     req.Rating,
		Review:     req.Review,
		RatingDate: req.RatingDate,
	})
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid book ID or reader ID.")
		}
		log.Error().Err(err).Msg("Failed to create rating")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create rating.", err.Error())
	}

	response := RatingResponse{
		ID:         rating.ID,
		BookID:     rating.BookID,
		ReaderID:   rating.ReaderID,
		Rating:     rating.Rating,
		Review:     rating.Review,
		RatingDate: rating.RatingDate,
		CreatedAt:  rating.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateRating обновляет рейтинг
func (h *Handler) updateRating(c *fiber.Ctx) error {
	idStr := c.Params("id")
	ratingID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid rating ID format.")
	}

	var req UpdateRatingRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию о рейтинге
	existingRating, err := h.repo.GetBookRatingByID(c.Context(), ratingID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Rating not found.")
		}
		log.Error().Err(err).Str("ratingID", idStr).Msg("Failed to get rating for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve rating.", err.Error())
	}

	// Проверка доступа: только владелец рейтинга может его изменить
	userIDStr := c.Locals("userID").(string)
	userID, _ := uuid.Parse(userIDStr)

	reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
	if err != nil || reader.ID != existingRating.ReaderID {
		return httperr.New(fiber.StatusForbidden, "Access denied. You can only update your own ratings.")
	}

	// Валидация
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		return httperr.New(fiber.StatusBadRequest, "Rating must be between 1 and 5.")
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateBookRatingParams{
		RatingID: ratingID,
	}

	if req.Rating != nil {
		updateParams.Rating = *req.Rating
	} else {
		updateParams.Rating = existingRating.Rating
	}

	if req.Review != nil {
		updateParams.Review = req.Review
	} else {
		updateParams.Review = existingRating.Review
	}

	// Обновление рейтинга
	updatedRating, err := h.repo.UpdateBookRating(c.Context(), updateParams)
	if err != nil {
		log.Error().Err(err).Str("ratingID", idStr).Msg("Failed to update rating")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update rating.", err.Error())
	}

	response := RatingResponse{
		ID:         updatedRating.ID,
		BookID:     updatedRating.BookID,
		ReaderID:   updatedRating.ReaderID,
		Rating:     updatedRating.Rating,
		Review:     updatedRating.Review,
		RatingDate: updatedRating.RatingDate,
		CreatedAt:  updatedRating.CreatedAt,
	}

	return c.JSON(response)
}

// deleteRating удаляет рейтинг
func (h *Handler) deleteRating(c *fiber.Ctx) error {
	idStr := c.Params("id")
	ratingID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid rating ID format.")
	}

	// Получаем информацию о рейтинге для проверки доступа
	existingRating, err := h.repo.GetBookRatingByID(c.Context(), ratingID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Rating not found.")
		}
		log.Error().Err(err).Str("ratingID", idStr).Msg("Failed to get rating for deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve rating.", err.Error())
	}

	// Проверка доступа
	userRole := c.Locals("userRole").(string)
	userIDStr := c.Locals("userID").(string)
	userID, _ := uuid.Parse(userIDStr)

	// Читатели могут удалять только свои рейтинги
	if userRole == string(postgres.UserRoleReader) {
		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != existingRating.ReaderID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only delete your own ratings.")
		}
	}
	// Библиотекари и администраторы могут удалять любые рейтинги

	// Удаление рейтинга
	err = h.repo.DeleteBookRating(c.Context(), ratingID)
	if err != nil {
		log.Error().Err(err).Str("ratingID", idStr).Msg("Failed to delete rating")
		return httperr.New(fiber.StatusInternalServerError, "Failed to delete rating.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

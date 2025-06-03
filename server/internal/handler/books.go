package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type CreateBookRequest struct {
	Title           string  `json:"title"`
	Author          string  `json:"author"`
	PublicationYear int     `json:"publication_year"`
	ISBN            *string `json:"isbn"`
	BookCode        string  `json:"book_code"`
	CategoryID      *int    `json:"category_id"`
	HallID          int     `json:"hall_id"`
	TotalCopies     int     `json:"total_copies"`
	ConditionStatus string  `json:"condition_status"`
	LocationInfo    *string `json:"location_info"`
}

type UpdateBookAvailabilityRequest struct {
	BookID          int64 `json:"book_id"`
	AvailableCopies int   `json:"available_copies"`
	PopularityScore int   `json:"popularity_score"`
}

type UpdateBookCopiesRequest struct {
	BookID          int64 `json:"book_id"`
	TotalCopies     int   `json:"total_copies"`
	AvailableCopies int   `json:"available_copies"`
}

type SearchBooksRequest struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	BookCode   string `json:"book_code"`
	ISBN       string `json:"isbn"`
	CategoryID int    `json:"category_id"`
	HallID     int    `json:"hall_id"`
	PageOffset int32  `json:"page_offset"`
	PageLimit  int32  `json:"page_limit"`
}

type AdvancedSearchBooksRequest struct {
	TitleFilter    string `json:"title_filter"`
	AuthorFilter   string `json:"author_filter"`
	YearFilter     int    `json:"year_filter"`
	CategoryFilter int    `json:"category_filter"`
	HallFilter     int    `json:"hall_filter"`
	AvailableOnly  bool   `json:"available_only"`
	SortBy         string `json:"sort_by"`
	PageOffset     int32  `json:"page_offset"`
	PageLimit      int32  `json:"page_limit"`
}

// createBook создает новую книгу в библиотеке
func (h *Handler) createBook(c *fiber.Ctx) error {
	req := new(CreateBookRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: title, author, book_code, hall_id, total_copies
	// TODO: Validate publication_year > 0 and <= current year
	// TODO: Validate total_copies > 0
	// TODO: Validate condition_status is one of: excellent, good, fair, poor, damaged
	// TODO: Validate ISBN format if provided
	// TODO: Validate hall_id exists
	// TODO: Validate category_id exists if provided

	params := postgres.CreateBookParams{
		Title:           req.Title,
		Author:          req.Author,
		PublicationYear: req.PublicationYear,
		Isbn:            req.ISBN,
		BookCode:        req.BookCode,
		CategoryID:      req.CategoryID,
		HallID:          req.HallID,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies, // Изначально все экземпляры доступны
		ConditionStatus: req.ConditionStatus,
		LocationInfo:    req.LocationInfo,
	}

	book, err := h.repo.CreateBook(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create book")
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Book with this code already exists.")
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid hall or category ID.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create book.")
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

// getBookByID получает книгу по ID
func (h *Handler) getBookByID(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	book, err := h.repo.GetBookByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Int64("bookID", id).Msg("Failed to get book by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.")
	}

	return c.JSON(book)
}

// getBookByCode получает книгу по коду
func (h *Handler) getBookByCode(c *fiber.Ctx) error {
	bookCode := c.Params("code")
	if bookCode == "" {
		return httperr.New(fiber.StatusBadRequest, "Book code is required.")
	}

	// TODO: Validate book_code format

	book, err := h.repo.GetBookByCode(c.Context(), bookCode)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookCode", bookCode).Msg("Failed to get book by code")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.")
	}

	return c.JSON(book)
}

// searchBooks выполняет поиск книг
func (h *Handler) searchBooks(c *fiber.Ctx) error {
	req := new(SearchBooksRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate at least one search parameter is provided
	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0

	if req.PageLimit == 0 {
		req.PageLimit = 20 // default limit
	}

	params := postgres.SearchBooksParams{
		Title:      req.Title,
		Author:     req.Author,
		BookCode:   req.BookCode,
		Isbn:       req.ISBN,
		CategoryID: req.CategoryID,
		HallID:     req.HallID,
		PageOffset: req.PageOffset,
		PageLimit:  req.PageLimit,
	}

	books, err := h.repo.SearchBooks(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Interface("params", params).Msg("Failed to search books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search books.")
	}

	if books == nil {
		books = []*postgres.SearchBooksRow{}
	}

	return c.JSON(books)
}

// advancedSearchBooks выполняет расширенный поиск книг
func (h *Handler) advancedSearchBooks(c *fiber.Ctx) error {
	req := new(AdvancedSearchBooksRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0
	// TODO: Validate year_filter if provided
	// TODO: Validate sort_by is one of: title, author, year, popularity
	// TODO: Validate hall_filter and category_filter if provided

	if req.PageLimit == 0 {
		req.PageLimit = 20 // default limit
	}

	if req.SortBy == "" {
		req.SortBy = "title" // default sort
	}

	params := postgres.AdvancedSearchBooksParams{
		TitleFilter:    req.TitleFilter,
		AuthorFilter:   req.AuthorFilter,
		YearFilter:     req.YearFilter,
		CategoryFilter: req.CategoryFilter,
		HallFilter:     req.HallFilter,
		AvailableOnly:  req.AvailableOnly,
		SortBy:         req.SortBy,
		PageOffset:     req.PageOffset,
		PageLimit:      req.PageLimit,
	}

	books, err := h.repo.AdvancedSearchBooks(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Interface("params", params).Msg("Failed to perform advanced search")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search books.")
	}

	if books == nil {
		books = []*postgres.AdvancedSearchBooksRow{}
	}

	return c.JSON(books)
}

// getAvailableBooks получает список доступных книг
func (h *Handler) getAvailableBooks(c *fiber.Ctx) error {
	limit := int32(20) // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// TODO: Validate limit > 0 and <= 100

	books, err := h.repo.GetAvailableBooks(c.Context(), limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get available books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve available books.")
	}

	if books == nil {
		books = []*postgres.GetAvailableBooksRow{}
	}

	return c.JSON(books)
}

// getPopularBooks получает список популярных книг
func (h *Handler) getPopularBooks(c *fiber.Ctx) error {
	limit := int32(10) // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// TODO: Validate limit > 0 and <= 50

	books, err := h.repo.GetPopularBooks(c.Context(), limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get popular books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve popular books.")
	}

	if books == nil {
		books = []*postgres.GetPopularBooksRow{}
	}

	return c.JSON(books)
}

// getTopRatedBooks получает список книг с высоким рейтингом
func (h *Handler) getTopRatedBooks(c *fiber.Ctx) error {
	limit := int32(10) // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// TODO: Validate limit > 0 and <= 50

	books, err := h.repo.GetTopRatedBooks(c.Context(), limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top rated books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve top rated books.")
	}

	if books == nil {
		books = []*postgres.GetTopRatedBooksRow{}
	}

	return c.JSON(books)
}

// updateBookAvailability обновляет количество доступных экземпляров книги
func (h *Handler) updateBookAvailability(c *fiber.Ctx) error {
	req := new(UpdateBookAvailabilityRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate book_id > 0
	// TODO: Validate available_copies >= 0
	// TODO: Validate popularity_score >= 0
	// TODO: Validate available_copies <= total_copies for the book

	params := postgres.UpdateBookAvailabilityParams{
		BookID:          req.BookID,
		AvailableCopies: req.AvailableCopies,
		PopularityScore: req.PopularityScore,
	}

	err := h.repo.UpdateBookAvailability(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Int64("bookID", req.BookID).Msg("Failed to update book availability")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book availability.")
	}

	return c.JSON(fiber.Map{"message": "Book availability updated successfully"})
}

// updateBookCopies обновляет общее количество экземпляров книги
func (h *Handler) updateBookCopies(c *fiber.Ctx) error {
	req := new(UpdateBookCopiesRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate book_id > 0
	// TODO: Validate total_copies > 0
	// TODO: Validate available_copies >= 0
	// TODO: Validate available_copies <= total_copies

	params := postgres.UpdateBookCopiesParams{
		BookID:          req.BookID,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.AvailableCopies,
	}

	err := h.repo.UpdateBookCopies(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Int64("bookID", req.BookID).Msg("Failed to update book copies")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book copies.")
	}

	return c.JSON(fiber.Map{"message": "Book copies updated successfully"})
}

// writeOffBook списывает книгу (уменьшает количество экземпляров)
func (h *Handler) writeOffBook(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.WriteOffBook(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Book not found or no copies available to write off.")
		}
		log.Error().Err(err).Int64("bookID", id).Msg("Failed to write off book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to write off book.")
	}

	return c.JSON(fiber.Map{"message": "Book successfully written off"})
}

// getBooksWithSingleCopy получает список читателей, взявших книги в единственном экземпляре
func (h *Handler) getBooksWithSingleCopy(c *fiber.Ctx) error {
	readers, err := h.repo.GetBooksWithSingleCopy(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get readers with single copy books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve readers with single copy books.")
	}

	if readers == nil {
		readers = []*postgres.GetBooksWithSingleCopyRow{}
	}

	return c.JSON(readers)
}

// getBooksByAuthorInHall получает статистику книг заданного автора в читальном зале
func (h *Handler) getBooksByAuthorInHall(c *fiber.Ctx) error {
	author := c.Query("author")
	hallIDStr := c.Query("hall_id")

	if author == "" || hallIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Author and hall_id are required parameters.")
	}

	// TODO: Validate author name
	// TODO: Validate hall_id is valid integer > 0

	hallID, err := strconv.Atoi(hallIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid hall_id parameter.")
	}

	params := postgres.GetBooksByAuthorInHallParams{
		Author: &author,
		HallID: hallID,
	}

	result, err := h.repo.GetBooksByAuthorInHall(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Str("author", author).Int("hallID", hallID).Msg("Failed to get books by author in hall")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve books count.")
	}

	return c.JSON(result)
}

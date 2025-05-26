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
	Title           string `json:"title"`            // validate: required, min=1, max=200
	Author          string `json:"author"`           // validate: required, min=1, max=100
	PublicationYear int32  `json:"publication_year"` // validate: required, min=1000, max=current_year
	Code            string `json:"code"`             // validate: required, unique, min=1, max=50
	CategoryID      int32  `json:"category_id"`      // validate: required, exists in categories
	TotalCopies     int32  `json:"total_copies"`     // validate: required, min=1, max=1000
	HallID          int32  `json:"hall_id"`          // validate: required, exists in halls
}

type UpdateBookAvailabilityRequest struct {
	AvailableCopies int32 `json:"available_copies"` // validate: required, min=0, max=total_copies
}

type SearchBooksRequest struct {
	Title  string `json:"title"`  // validate: optional, max=200
	Author string `json:"author"` // validate: optional, max=100
	Limit  int32  `json:"limit"`  // validate: optional, min=1, max=100, default=20
}

func (h *Handler) createBook(c *fiber.Ctx) error {
	req := new(CreateBookRequest)
	// TODO: validate required fields (title, author, publication_year, code, category_id, total_copies, hall_id)
	// TODO: validate publication_year range (1000 - current year)
	// TODO: validate code uniqueness
	// TODO: validate category_id exists in book_categories table
	// TODO: validate hall_id exists in halls table
	// TODO: validate total_copies range (1-1000)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	params := postgres.CreateBookParams{
		Title:           req.Title,
		Author:          req.Author,
		PublicationYear: req.PublicationYear,
		BookCode:        req.Code,
		CategoryID:      req.CategoryID,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies, // Initially all copies are available
		HallID:          req.HallID,
	}

	book, err := h.repo.CreateBook(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create book")
		if strings.Contains(err.Error(), "unique constraint") && strings.Contains(err.Error(), "code") {
			return httperr.New(fiber.StatusConflict, "Book with this code already exists.")
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid category or hall ID.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create book.")
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func (h *Handler) getBookByID(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	book, err := h.repo.GetBookByID(c.Context(), int32(id))
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Int64("bookID", id).Msg("Failed to get book by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.")
	}

	return c.JSON(book)
}

func (h *Handler) getBookByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	// TODO: validate code format and length (min=1, max=50)
	if code == "" {
		return httperr.New(fiber.StatusBadRequest, "Book code is required.")
	}

	book, err := h.repo.GetBookByCode(c.Context(), code)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookCode", code).Msg("Failed to get book by code")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.")
	}

	return c.JSON(book)
}

func (h *Handler) searchBooks(c *fiber.Ctx) error {
	req := new(SearchBooksRequest)
	// TODO: validate title and author length (max=200 and max=100 respectively)
	// TODO: validate limit range (1-100), set default to 20
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	if req.Limit == 0 {
		req.Limit = 20 // Default limit
	}

	params := postgres.SearchBooksParams{
		Title:     req.Title,
		Author:    req.Author,
		PageLimit: req.Limit,
	}

	books, err := h.repo.SearchBooks(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search books.")
	}

	if books == nil {
		books = []*postgres.SearchBooksRow{}
	}

	return c.JSON(books)
}

func (h *Handler) getAvailableBooks(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "50")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit <= 0 {
		limit = 50 // Default limit
	}
	// TODO: validate limit range (1-100)

	books, err := h.repo.GetAvailableBooks(c.Context(), int32(limit))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get available books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve available books.")
	}

	if books == nil {
		books = []*postgres.GetAvailableBooksRow{}
	}

	return c.JSON(books)
}

func (h *Handler) getPopularBooks(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}
	// TODO: validate limit range (1-50)

	books, err := h.repo.GetPopularBooks(c.Context(), int32(limit))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get popular books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve popular books.")
	}

	if books == nil {
		books = []*postgres.GetPopularBooksRow{}
	}

	return c.JSON(books)
}

func (h *Handler) updateBookAvailability(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(UpdateBookAvailabilityRequest)
	// TODO: validate available_copies range (0 - total_copies of the book)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	params := postgres.UpdateBookAvailabilityParams{
		BookID:          int32(id),
		AvailableCopies: req.AvailableCopies,
	}

	err = h.repo.UpdateBookAvailability(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Int64("bookID", id).Msg("Failed to update book availability")
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book availability.")
	}

	return c.JSON(fiber.Map{
		"message":          "Book availability updated successfully",
		"book_id":          id,
		"available_copies": req.AvailableCopies,
	})
}

func (h *Handler) globalSearch(c *fiber.Ctx) error {
	searchTerm := c.Query("q")
	// TODO: validate search term (required, min=1, max=100, sanitize for SQL injection)
	if searchTerm == "" {
		return httperr.New(fiber.StatusBadRequest, "Search term is required.")
	}

	results, err := h.repo.GlobalSearch(c.Context(), searchTerm)
	if err != nil {
		log.Error().Err(err).Str("searchTerm", searchTerm).Msg("Failed to perform global search")
		return httperr.New(fiber.StatusInternalServerError, "Failed to perform search.")
	}

	if results == nil {
		results = []*postgres.GlobalSearchRow{}
	}

	return c.JSON(results)
}

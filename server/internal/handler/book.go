package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

type CreateBookRequest struct {
	Title           string   `json:"title" validate:"required"`
	ISBN            *string  `json:"isbn"`
	PublicationYear *int     `json:"publication_year"`
	Publisher       *string  `json:"publisher"`
	TotalCopies     int      `json:"total_copies" validate:"required,min=1"`
	Authors         []string `json:"authors"`
}

type UpdateBookRequest struct {
	Title           string  `json:"title" validate:"required"`
	ISBN            *string `json:"isbn"`
	PublicationYear *int    `json:"publication_year"`
	Publisher       *string `json:"publisher"`
}

type BookWithAuthorsResponse struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	ISBN            *string   `json:"isbn"`
	PublicationYear *int      `json:"publication_year"`
	Publisher       *string   `json:"publisher"`
	TotalCopies     int       `json:"total_copies"`
	AvailableCopies int       `json:"available_copies"`
	Authors         string    `json:"authors"`
}

func (h *Handler) getAllBooks(c *fiber.Ctx) error {
	books, err := h.repo.GetAllBooks(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve books")
	}

	response := make([]BookWithAuthorsResponse, len(books))
	for i, book := range books {
		response[i] = BookWithAuthorsResponse{
			ID:              book.ID,
			Title:           book.Title,
			ISBN:            book.Isbn,
			PublicationYear: book.PublicationYear,
			Publisher:       book.Publisher,
			TotalCopies:     book.TotalCopies,
			AvailableCopies: book.AvailableCopies,
			Authors:         string(book.Authors),
		}
	}

	return c.JSON(response)
}

func (h *Handler) searchBooks(c *fiber.Ctx) error {
	title := c.Query("title")
	author := c.Query("author")
	yearStr := c.Query("year")

	var year int
	if yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid year parameter")
		}
	}

	books, err := h.repo.SearchBooks(c.Context(), postgres.SearchBooksParams{
		Title:           title,
		Author:          author,
		PublicationYear: year,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to search books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search books")
	}

	response := make([]BookWithAuthorsResponse, len(books))
	for i, book := range books {
		response[i] = BookWithAuthorsResponse{
			ID:              book.ID,
			Title:           book.Title,
			ISBN:            book.Isbn,
			PublicationYear: book.PublicationYear,
			Publisher:       book.Publisher,
			TotalCopies:     book.TotalCopies,
			AvailableCopies: book.AvailableCopies,
			Authors:         string(book.Authors),
		}
	}

	return c.JSON(response)
}

func (h *Handler) getBookById(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format")
	}

	book, err := h.repo.GetBookById(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book")
	}

	return c.JSON(book)
}

func (h *Handler) createBook(c *fiber.Ctx) error {
	var req CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Create book
	book, err := h.repo.CreateBook(c.Context(), postgres.CreateBookParams{
		Title:           req.Title,
		Isbn:            req.ISBN,
		PublicationYear: req.PublicationYear,
		Publisher:       req.Publisher,
		TotalCopies:     req.TotalCopies,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create book")
	}

	// Add authors if provided
	for _, authorName := range req.Authors {
		if authorName == "" {
			continue
		}

		author, err := h.repo.GetOrCreateAuthor(c.Context(), authorName)
		if err != nil {
			log.Warn().Err(err).Str("authorName", authorName).Msg("Failed to create author")
			continue
		}

		err = h.repo.AddBookAuthor(c.Context(), postgres.AddBookAuthorParams{
			BookID:   book.ID,
			AuthorID: author.ID,
		})
		if err != nil {
			log.Warn().Err(err).Str("authorID", author.ID.String()).Msg("Failed to add book author")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func (h *Handler) updateBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format")
	}

	var req UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	book, err := h.repo.UpdateBook(c.Context(), postgres.UpdateBookParams{
		ID:              id,
		Title:           req.Title,
		Isbn:            req.ISBN,
		PublicationYear: req.PublicationYear,
		Publisher:       req.Publisher,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to update book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book")
	}

	return c.JSON(book)
}

func (h *Handler) getBookAuthors(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format")
	}

	authors, err := h.repo.GetBookAuthors(c.Context(), id)
	if err != nil {
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get book authors")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book authors")
	}

	return c.JSON(authors)
}

type AddBookAuthorRequest struct {
	AuthorID string `json:"author_id" validate:"required"`
}

func (h *Handler) addBookAuthor(c *fiber.Ctx) error {
	bookIdStr := c.Params("id")
	bookId, err := uuid.Parse(bookIdStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format")
	}

	var req AddBookAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	authorId, err := uuid.Parse(req.AuthorID)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format")
	}

	err = h.repo.AddBookAuthor(c.Context(), postgres.AddBookAuthorParams{
		BookID:   bookId,
		AuthorID: authorId,
	})
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIdStr).Str("authorID", req.AuthorID).Msg("Failed to add book author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to add book author")
	}

	return c.JSON(fiber.Map{"message": "Author added to book successfully"})
}

func (h *Handler) removeBookAuthor(c *fiber.Ctx) error {
	bookIdStr := c.Params("id")
	bookId, err := uuid.Parse(bookIdStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format")
	}

	authorIdStr := c.Params("authorId")
	authorId, err := uuid.Parse(authorIdStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format")
	}

	err = h.repo.RemoveBookAuthor(c.Context(), postgres.RemoveBookAuthorParams{
		BookID:   bookId,
		AuthorID: authorId,
	})
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIdStr).Str("authorID", authorIdStr).Msg("Failed to remove book author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to remove book author")
	}

	return c.JSON(fiber.Map{"message": "Author removed from book successfully"})
}

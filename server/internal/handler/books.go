package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

// Request/Response structs для books
type CreateBookRequest struct {
	Title           string      `json:"title" validate:"required,min=1,max=500"`
	ISBN            *string     `json:"isbn" validate:"omitempty,len=13"`
	PublicationYear *int        `json:"publication_year" validate:"omitempty,min=1000,max=2030"`
	Publisher       *string     `json:"publisher" validate:"omitempty,max=200"`
	Pages           *int        `json:"pages" validate:"omitempty,min=1,max=10000"`
	Language        *string     `json:"language" validate:"omitempty,max=50"`
	Description     *string     `json:"description" validate:"omitempty,max=2000"`
	TotalCopies     int         `json:"total_copies" validate:"required,min=1,max=1000"`
	AvailableCopies int         `json:"available_copies" validate:"required,min=0"`
	AuthorIDs       []uuid.UUID `json:"author_ids" validate:"omitempty"`
}

type UpdateBookRequest struct {
	Title           *string     `json:"title" validate:"omitempty,min=1,max=500"`
	ISBN            *string     `json:"isbn" validate:"omitempty,len=13"`
	PublicationYear *int        `json:"publication_year" validate:"omitempty,min=1000,max=2030"`
	Publisher       *string     `json:"publisher" validate:"omitempty,max=200"`
	Pages           *int        `json:"pages" validate:"omitempty,min=1,max=10000"`
	Language        *string     `json:"language" validate:"omitempty,max=50"`
	Description     *string     `json:"description" validate:"omitempty,max=2000"`
	TotalCopies     *int        `json:"total_copies" validate:"omitempty,min=1,max=1000"`
	AvailableCopies *int        `json:"available_copies" validate:"omitempty,min=0"`
	AuthorIDs       []uuid.UUID `json:"author_ids" validate:"omitempty"`
}

type UpdateBookAvailabilityRequest struct {
	AvailableCopies int `json:"available_copies" validate:"required,min=0"`
}

type BookResponse struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	ISBN            *string    `json:"isbn"`
	PublicationYear *int       `json:"publication_year"`
	Publisher       *string    `json:"publisher"`
	Pages           *int       `json:"pages"`
	Language        *string    `json:"language"`
	Description     *string    `json:"description"`
	TotalCopies     int        `json:"total_copies"`
	AvailableCopies int        `json:"available_copies"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	// Extended fields
	Authors     *string          `json:"authors,omitempty"`
	AvgRating   *decimal.Decimal `json:"avg_rating,omitempty"`
	RatingCount *int64           `json:"rating_count,omitempty"`
}

type BooksListResponse struct {
	Books  []BookResponse `json:"books"`
	Total  int64          `json:"total"`
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
}

// listBooks возвращает список книг с пагинацией
func (h *Handler) listBooks(c *fiber.Ctx) error {
	// Парсинг параметров пагинации
	limit := int32(20)
	offset := int32(0)
	withAuthors := c.Query("with_authors", "false") == "true"

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

	var books []BookResponse
	var err error

	if withAuthors {
		// Получение книг с авторами
		booksWithAuthors, err := h.repo.GetBooksWithAuthors(c.Context(), postgres.GetBooksWithAuthorsParams{
			LimitVal:  limit,
			OffsetVal: offset,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to list books with authors")
			return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve books.", err.Error())
		}

		books = make([]BookResponse, len(booksWithAuthors))
		for i, book := range booksWithAuthors {
			authors := string(book.Authors)
			if authors == "" {
				authors = "Unknown Author"
			}

			books[i] = BookResponse{
				ID:              book.ID,
				Title:           book.Title,
				ISBN:            book.Isbn,
				PublicationYear: book.PublicationYear,
				Publisher:       book.Publisher,
				Pages:           book.Pages,
				Language:        book.Language,
				Description:     book.Description,
				TotalCopies:     book.TotalCopies,
				AvailableCopies: book.AvailableCopies,
				CreatedAt:       book.CreatedAt,
				UpdatedAt:       book.UpdatedAt,
				Authors:         &authors,
			}
		}
	} else {
		// Получение обычных книг
		simpleBooks, err := h.repo.ListBooks(c.Context(), postgres.ListBooksParams{
			LimitVal:  limit,
			OffsetVal: offset,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to list books")
			return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve books.", err.Error())
		}

		books = make([]BookResponse, len(simpleBooks))
		for i, book := range simpleBooks {
			books[i] = BookResponse{
				ID:              book.ID,
				Title:           book.Title,
				ISBN:            book.Isbn,
				PublicationYear: book.PublicationYear,
				Publisher:       book.Publisher,
				Pages:           book.Pages,
				Language:        book.Language,
				Description:     book.Description,
				TotalCopies:     book.TotalCopies,
				AvailableCopies: book.AvailableCopies,
				CreatedAt:       book.CreatedAt,
				UpdatedAt:       book.UpdatedAt,
			}
		}
	}

	// Получение общего количества книг
	total, err := h.repo.CountBooks(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to count books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to count books.", err.Error())
	}

	response := BooksListResponse{
		Books:  books,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	return c.JSON(response)
}

// searchBooks ищет книги по названию
func (h *Handler) searchBooks(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return httperr.New(fiber.StatusBadRequest, "Search query parameter 'q' is required.")
	}

	books, err := h.repo.SearchBooksByTitle(c.Context(), query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to search books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search books.", err.Error())
	}

	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = BookResponse{
			ID:              book.ID,
			Title:           book.Title,
			ISBN:            book.Isbn,
			PublicationYear: book.PublicationYear,
			Publisher:       book.Publisher,
			Pages:           book.Pages,
			Language:        book.Language,
			Description:     book.Description,
			TotalCopies:     book.TotalCopies,
			AvailableCopies: book.AvailableCopies,
			CreatedAt:       book.CreatedAt,
			UpdatedAt:       book.UpdatedAt,
		}
	}

	return c.JSON(fiber.Map{"books": response})
}

// getTopRatedBooks возвращает топ книг по рейтингу
func (h *Handler) getTopRatedBooks(c *fiber.Ctx) error {
	limit := int32(10)
	minRatings := 1

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = int32(parsedLimit)
		}
	}

	if mr := c.Query("min_ratings"); mr != "" {
		if parsedMinRatings, err := strconv.Atoi(mr); err == nil && parsedMinRatings > 0 {
			minRatings = parsedMinRatings
		}
	}

	books, err := h.repo.GetTopRatedBooks(c.Context(), postgres.GetTopRatedBooksParams{
		MinRatings: minRatings,
		LimitVal:   limit,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top rated books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve top rated books.", err.Error())
	}

	response := make([]BookResponse, len(books))
	for i, book := range books {
		authors := string(book.Authors)
		if authors == "" {
			authors = "Unknown Author"
		}

		response[i] = BookResponse{
			ID:              book.ID,
			Title:           book.Title,
			ISBN:            book.Isbn,
			PublicationYear: book.PublicationYear,
			Publisher:       book.Publisher,
			Pages:           book.Pages,
			Language:        book.Language,
			Description:     book.Description,
			TotalCopies:     book.TotalCopies,
			AvailableCopies: book.AvailableCopies,
			CreatedAt:       book.CreatedAt,
			UpdatedAt:       book.UpdatedAt,
			Authors:         &authors,
			AvgRating:       &book.AvgRating,
			RatingCount:     &book.RatingCount,
		}
	}

	return c.JSON(fiber.Map{"books": response})
}

// getBooksByAuthor возвращает книги автора
func (h *Handler) getBooksByAuthor(c *fiber.Ctx) error {
	authorIDStr := c.Params("authorId")
	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format.")
	}

	books, err := h.repo.GetBooksByAuthor(c.Context(), authorID)
	if err != nil {
		log.Error().Err(err).Str("authorID", authorIDStr).Msg("Failed to get books by author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve books by author.", err.Error())
	}

	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = BookResponse{
			ID:              book.ID,
			Title:           book.Title,
			ISBN:            book.Isbn,
			PublicationYear: book.PublicationYear,
			Publisher:       book.Publisher,
			Pages:           book.Pages,
			Language:        book.Language,
			Description:     book.Description,
			TotalCopies:     book.TotalCopies,
			AvailableCopies: book.AvailableCopies,
			CreatedAt:       book.CreatedAt,
			UpdatedAt:       book.UpdatedAt,
		}
	}

	return c.JSON(fiber.Map{"books": response})
}

// getBookByISBN возвращает книгу по ISBN
func (h *Handler) getBookByISBN(c *fiber.Ctx) error {
	isbn := c.Params("isbn")
	if isbn == "" {
		return httperr.New(fiber.StatusBadRequest, "ISBN is required.")
	}

	book, err := h.repo.GetBookByISBN(c.Context(), &isbn)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book with this ISBN not found.")
		}
		log.Error().Err(err).Str("isbn", isbn).Msg("Failed to get book by ISBN")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.", err.Error())
	}

	response := BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		ISBN:            book.Isbn,
		PublicationYear: book.PublicationYear,
		Publisher:       book.Publisher,
		Pages:           book.Pages,
		Language:        book.Language,
		Description:     book.Description,
		TotalCopies:     book.TotalCopies,
		AvailableCopies: book.AvailableCopies,
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
	}

	return c.JSON(response)
}

// getBookByID возвращает книгу по ID
func (h *Handler) getBookByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	bookID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	book, err := h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get book by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.", err.Error())
	}

	response := BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		ISBN:            book.Isbn,
		PublicationYear: book.PublicationYear,
		Publisher:       book.Publisher,
		Pages:           book.Pages,
		Language:        book.Language,
		Description:     book.Description,
		TotalCopies:     book.TotalCopies,
		AvailableCopies: book.AvailableCopies,
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
	}

	return c.JSON(response)
}

// getBookWithDetails возвращает книгу с деталями (авторы, рейтинг)
func (h *Handler) getBookWithDetails(c *fiber.Ctx) error {
	idStr := c.Params("id")
	bookID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	book, err := h.repo.GetBookWithDetails(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get book with details")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book details.", err.Error())
	}

	authors := string(book.Authors)
	if authors == "" {
		authors = "Unknown Author"
	}

	// Конвертируем interface{} в decimal.Decimal для рейтинга
	var avgRating *decimal.Decimal
	if book.AvgRating != nil {
		if rating, err := decimal.Parse(book.AvgRating.(string)); err == nil {
			avgRating = &rating
		}
	}

	response := BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		ISBN:            book.Isbn,
		PublicationYear: book.PublicationYear,
		Publisher:       book.Publisher,
		Pages:           book.Pages,
		Language:        book.Language,
		Description:     book.Description,
		TotalCopies:     book.TotalCopies,
		AvailableCopies: book.AvailableCopies,
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
		Authors:         &authors,
		AvgRating:       avgRating,
		RatingCount:     &book.RatingCount,
	}

	return c.JSON(response)
}

// createBook создает новую книгу
func (h *Handler) createBook(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can create books.")
	}

	var req CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация
	if req.Title == "" {
		return httperr.New(fiber.StatusBadRequest, "Book title is required.")
	}
	if req.TotalCopies <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Total copies must be greater than 0.")
	}
	if req.AvailableCopies < 0 || req.AvailableCopies > req.TotalCopies {
		return httperr.New(fiber.StatusBadRequest, "Available copies must be between 0 and total copies.")
	}

	// Создание книги
	book, err := h.repo.CreateBook(c.Context(), postgres.CreateBookParams{
		Title:           req.Title,
		Isbn:            req.ISBN,
		PublicationYear: req.PublicationYear,
		Publisher:       req.Publisher,
		Pages:           req.Pages,
		Language:        req.Language,
		Description:     req.Description,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.AvailableCopies,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "isbn") {
			return httperr.New(fiber.StatusConflict, "Book with this ISBN already exists.")
		}
		log.Error().Err(err).Msg("Failed to create book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create book.", err.Error())
	}

	// Добавление авторов, если указаны
	if len(req.AuthorIDs) > 0 {
		for _, authorID := range req.AuthorIDs {
			err := h.repo.AddBookAuthor(c.Context(), postgres.AddBookAuthorParams{
				BookID:   book.ID,
				AuthorID: authorID,
			})
			if err != nil {
				log.Warn().Err(err).Str("bookID", book.ID.String()).Str("authorID", authorID.String()).Msg("Failed to add author to book")
			}
		}
	}

	response := BookResponse{
		ID:              book.ID,
		Title:           book.Title,
		ISBN:            book.Isbn,
		PublicationYear: book.PublicationYear,
		Publisher:       book.Publisher,
		Pages:           book.Pages,
		Language:        book.Language,
		Description:     book.Description,
		TotalCopies:     book.TotalCopies,
		AvailableCopies: book.AvailableCopies,
		CreatedAt:       book.CreatedAt,
		UpdatedAt:       book.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateBook обновляет информацию о книге
func (h *Handler) updateBook(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can update books.")
	}

	idStr := c.Params("id")
	bookID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	var req UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию о книге
	existingBook, err := h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get book for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.", err.Error())
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateBookParams{
		BookID: bookID,
	}

	// Устанавливаем значения для обновления или сохраняем существующие
	if req.Title != nil {
		if *req.Title == "" {
			return httperr.New(fiber.StatusBadRequest, "Book title cannot be empty.")
		}
		updateParams.Title = *req.Title
	} else {
		updateParams.Title = existingBook.Title
	}

	if req.ISBN != nil {
		updateParams.Isbn = req.ISBN
	} else {
		updateParams.Isbn = existingBook.Isbn
	}

	if req.PublicationYear != nil {
		updateParams.PublicationYear = req.PublicationYear
	} else {
		updateParams.PublicationYear = existingBook.PublicationYear
	}

	if req.Publisher != nil {
		updateParams.Publisher = req.Publisher
	} else {
		updateParams.Publisher = existingBook.Publisher
	}

	if req.Pages != nil {
		updateParams.Pages = req.Pages
	} else {
		updateParams.Pages = existingBook.Pages
	}

	if req.Language != nil {
		updateParams.Language = req.Language
	} else {
		updateParams.Language = existingBook.Language
	}

	if req.Description != nil {
		updateParams.Description = req.Description
	} else {
		updateParams.Description = existingBook.Description
	}

	if req.TotalCopies != nil {
		if *req.TotalCopies <= 0 {
			return httperr.New(fiber.StatusBadRequest, "Total copies must be greater than 0.")
		}
		updateParams.TotalCopies = *req.TotalCopies
	} else {
		updateParams.TotalCopies = existingBook.TotalCopies
	}

	if req.AvailableCopies != nil {
		if *req.AvailableCopies < 0 || *req.AvailableCopies > updateParams.TotalCopies {
			return httperr.New(fiber.StatusBadRequest, "Available copies must be between 0 and total copies.")
		}
		updateParams.AvailableCopies = *req.AvailableCopies
	} else {
		updateParams.AvailableCopies = existingBook.AvailableCopies
	}

	// Обновление книги
	updatedBook, err := h.repo.UpdateBook(c.Context(), updateParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "isbn") {
			return httperr.New(fiber.StatusConflict, "Book with this ISBN already exists.")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to update book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book.", err.Error())
	}

	// Обновление авторов, если указаны
	if req.AuthorIDs != nil {
		// Удаляем всех текущих авторов
		err := h.repo.RemoveAllBookAuthors(c.Context(), bookID)
		if err != nil {
			log.Warn().Err(err).Str("bookID", idStr).Msg("Failed to remove book authors")
		}

		// Добавляем новых авторов
		for _, authorID := range req.AuthorIDs {
			err := h.repo.AddBookAuthor(c.Context(), postgres.AddBookAuthorParams{
				BookID:   bookID,
				AuthorID: authorID,
			})
			if err != nil {
				log.Warn().Err(err).Str("bookID", idStr).Str("authorID", authorID.String()).Msg("Failed to add author to book")
			}
		}
	}

	response := BookResponse{
		ID:              updatedBook.ID,
		Title:           updatedBook.Title,
		ISBN:            updatedBook.Isbn,
		PublicationYear: updatedBook.PublicationYear,
		Publisher:       updatedBook.Publisher,
		Pages:           updatedBook.Pages,
		Language:        updatedBook.Language,
		Description:     updatedBook.Description,
		TotalCopies:     updatedBook.TotalCopies,
		AvailableCopies: updatedBook.AvailableCopies,
		CreatedAt:       updatedBook.CreatedAt,
		UpdatedAt:       updatedBook.UpdatedAt,
	}

	return c.JSON(response)
}

// updateBookAvailability обновляет количество доступных экземпляров
func (h *Handler) updateBookAvailability(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can update book availability.")
	}

	idStr := c.Params("id")
	bookID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	var req UpdateBookAvailabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Проверяем, что книга существует и получаем информацию о ней
	existingBook, err := h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get book for availability update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.", err.Error())
	}

	// Валидация
	if req.AvailableCopies < 0 || req.AvailableCopies > existingBook.TotalCopies {
		return httperr.New(fiber.StatusBadRequest, "Available copies must be between 0 and total copies.")
	}

	// Обновление доступности
	err = h.repo.UpdateBookAvailability(c.Context(), postgres.UpdateBookAvailabilityParams{
		BookID:          bookID,
		AvailableCopies: req.AvailableCopies,
	})
	if err != nil {
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to update book availability")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book availability.", err.Error())
	}

	// Получаем обновленную информацию
	updatedBook, err := h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get updated book info")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve updated book info.", err.Error())
	}

	response := BookResponse{
		ID:              updatedBook.ID,
		Title:           updatedBook.Title,
		ISBN:            updatedBook.Isbn,
		PublicationYear: updatedBook.PublicationYear,
		Publisher:       updatedBook.Publisher,
		Pages:           updatedBook.Pages,
		Language:        updatedBook.Language,
		Description:     updatedBook.Description,
		TotalCopies:     updatedBook.TotalCopies,
		AvailableCopies: updatedBook.AvailableCopies,
		CreatedAt:       updatedBook.CreatedAt,
		UpdatedAt:       updatedBook.UpdatedAt,
	}

	return c.JSON(response)
}

// deleteBook удаляет книгу
func (h *Handler) deleteBook(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators can delete books.")
	}

	idStr := c.Params("id")
	bookID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	// Проверяем, что книга существует
	_, err = h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to get book for deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book.", err.Error())
	}

	// Удаляем связи с авторами
	err = h.repo.RemoveAllBookAuthors(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to remove book authors before deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to prepare book for deletion.", err.Error())
	}

	// Удаляем книгу
	err = h.repo.DeleteBook(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusConflict, "Cannot delete book that has active issues or copies. Please remove all dependencies first.")
		}
		log.Error().Err(err).Str("bookID", idStr).Msg("Failed to delete book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to delete book.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// Добавляем в books.go или в отдельный хендлер
// getTopRatedBooksWithRatings возвращает топ книг с рейтингами (более детальная информация)
func (h *Handler) getTopRatedBooksWithRatings(c *fiber.Ctx) error {
	limit := int32(10)
	minRatings := 1

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = int32(parsedLimit)
		}
	}

	if mr := c.Query("min_ratings"); mr != "" {
		if parsedMinRatings, err := strconv.Atoi(mr); err == nil && parsedMinRatings > 0 {
			minRatings = parsedMinRatings
		}
	}

	books, err := h.repo.GetTopRatedBooksWithRatings(c.Context(), postgres.GetTopRatedBooksWithRatingsParams{
		MinRatings: minRatings,
		LimitVal:   limit,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top rated books with ratings")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve top rated books.", err.Error())
	}

	response := make([]TopRatedBookResponse, len(books))
	for i, book := range books {
		authors := string(book.Authors)
		if authors == "" {
			authors = "Unknown Author"
		}

		response[i] = TopRatedBookResponse{
			ID:          book.ID,
			Title:       book.Title,
			Authors:     authors,
			AvgRating:   book.AvgRating,
			RatingCount: book.RatingCount,
		}
	}

	return c.JSON(fiber.Map{"books": response})
}

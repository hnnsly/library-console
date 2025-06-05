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

// Request/Response structs для authors
type CreateAuthorRequest struct {
	FullName  string  `json:"full_name" validate:"required,min=2,max=200"`
	BirthYear *int    `json:"birth_year" validate:"omitempty,min=1,max=2030"`
	DeathYear *int    `json:"death_year" validate:"omitempty,min=1,max=2030"`
	Biography *string `json:"biography" validate:"omitempty,max=5000"`
}

type UpdateAuthorRequest struct {
	FullName  *string `json:"full_name" validate:"omitempty,min=2,max=200"`
	BirthYear *int    `json:"birth_year" validate:"omitempty,min=1,max=2030"`
	DeathYear *int    `json:"death_year" validate:"omitempty,min=1,max=2030"`
	Biography *string `json:"biography" validate:"omitempty,max=5000"`
}

type AddBookAuthorRequest struct {
	AuthorID uuid.UUID `json:"author_id" validate:"required"`
}

type AuthorResponse struct {
	ID        uuid.UUID  `json:"id"`
	FullName  string     `json:"full_name"`
	BirthYear *int       `json:"birth_year"`
	DeathYear *int       `json:"death_year"`
	Biography *string    `json:"biography"`
	CreatedAt *time.Time `json:"created_at"`
	// Extended fields
	BooksCount *int64 `json:"books_count,omitempty"`
	Age        *int   `json:"age,omitempty"`
	IsAlive    bool   `json:"is_alive"`
}

type AuthorWithBooksResponse struct {
	AuthorResponse
	Books []BookResponse `json:"books"`
}

type AuthorsListResponse struct {
	Authors []AuthorResponse `json:"authors"`
	Total   int64            `json:"total"`
	Limit   int32            `json:"limit"`
	Offset  int32            `json:"offset"`
}

// calculateAge вычисляет возраст автора
func calculateAge(birthYear *int, deathYear *int) *int {
	if birthYear == nil {
		return nil
	}

	endYear := time.Now().Year()
	if deathYear != nil {
		endYear = *deathYear
	}

	age := endYear - *birthYear
	return &age
}

// listAuthors возвращает список авторов с пагинацией
func (h *Handler) listAuthors(c *fiber.Ctx) error {
	// Пагинация
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

	authors, err := h.repo.ListAuthors(c.Context(), postgres.ListAuthorsParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to list authors")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve authors.", err.Error())
	}

	// Получаем общее количество авторов
	total, err := h.repo.CountAuthors(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to count authors")
		total = 0 // Не критично, продолжаем
	}

	response := make([]AuthorResponse, len(authors))
	for i, author := range authors {
		age := calculateAge(author.BirthYear, author.DeathYear)
		isAlive := author.DeathYear == nil

		response[i] = AuthorResponse{
			ID:        author.ID,
			FullName:  author.FullName,
			BirthYear: author.BirthYear,
			DeathYear: author.DeathYear,
			Biography: author.Biography,
			CreatedAt: author.CreatedAt,
			Age:       age,
			IsAlive:   isAlive,
		}
	}

	return c.JSON(AuthorsListResponse{
		Authors: response,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// searchAuthors выполняет поиск авторов по имени
func (h *Handler) searchAuthors(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return httperr.New(fiber.StatusBadRequest, "Search query parameter 'q' is required.")
	}

	if len(query) < 2 {
		return httperr.New(fiber.StatusBadRequest, "Search query must be at least 2 characters long.")
	}

	authors, err := h.repo.SearchAuthorsByName(c.Context(), query)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to search authors")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search authors.", err.Error())
	}

	response := make([]AuthorResponse, len(authors))
	for i, author := range authors {
		age := calculateAge(author.BirthYear, author.DeathYear)
		isAlive := author.DeathYear == nil

		response[i] = AuthorResponse{
			ID:        author.ID,
			FullName:  author.FullName,
			BirthYear: author.BirthYear,
			DeathYear: author.DeathYear,
			Biography: author.Biography,
			CreatedAt: author.CreatedAt,
			Age:       age,
			IsAlive:   isAlive,
		}
	}

	return c.JSON(fiber.Map{
		"authors": response,
		"query":   query,
		"count":   len(response),
	})
}

// getAuthorByID возвращает автора по ID
func (h *Handler) getAuthorByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	authorID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format.")
	}

	author, err := h.repo.GetAuthorByID(c.Context(), authorID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Author not found.")
		}
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to get author by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve author.", err.Error())
	}

	// Получаем количество книг автора
	books, err := h.repo.GetAuthorBooks(c.Context(), authorID)
	if err != nil {
		log.Warn().Err(err).Str("authorID", idStr).Msg("Failed to get author books count")
	}

	var booksCount *int64
	if books != nil {
		count := int64(len(books))
		booksCount = &count
	}

	age := calculateAge(author.BirthYear, author.DeathYear)
	isAlive := author.DeathYear == nil

	response := AuthorResponse{
		ID:         author.ID,
		FullName:   author.FullName,
		BirthYear:  author.BirthYear,
		DeathYear:  author.DeathYear,
		Biography:  author.Biography,
		CreatedAt:  author.CreatedAt,
		BooksCount: booksCount,
		Age:        age,
		IsAlive:    isAlive,
	}

	return c.JSON(response)
}

// getAuthorBooks возвращает книги автора
func (h *Handler) getAuthorBooks(c *fiber.Ctx) error {
	idStr := c.Params("id")
	authorID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format.")
	}

	// Проверяем, что автор существует
	author, err := h.repo.GetAuthorByID(c.Context(), authorID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Author not found.")
		}
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to verify author existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify author.", err.Error())
	}

	books, err := h.repo.GetAuthorBooks(c.Context(), authorID)
	if err != nil {
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to get author books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve author books.", err.Error())
	}

	// Конвертируем в BookResponse (предполагаем, что структура BookResponse уже определена)
	bookResponses := make([]BookResponse, len(books))
	for i, book := range books {
		bookResponses[i] = BookResponse{
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

	age := calculateAge(author.BirthYear, author.DeathYear)
	isAlive := author.DeathYear == nil
	booksCount := int64(len(books))

	authorResponse := AuthorResponse{
		ID:         author.ID,
		FullName:   author.FullName,
		BirthYear:  author.BirthYear,
		DeathYear:  author.DeathYear,
		Biography:  author.Biography,
		CreatedAt:  author.CreatedAt,
		BooksCount: &booksCount,
		Age:        age,
		IsAlive:    isAlive,
	}

	response := AuthorWithBooksResponse{
		AuthorResponse: authorResponse,
		Books:          bookResponses,
	}

	return c.JSON(response)
}

// createAuthor создает нового автора
func (h *Handler) createAuthor(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can create authors.")
	}

	var req CreateAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация
	if req.FullName == "" {
		return httperr.New(fiber.StatusBadRequest, "Author full name is required.")
	}

	if len(req.FullName) < 2 {
		return httperr.New(fiber.StatusBadRequest, "Author full name must be at least 2 characters long.")
	}

	// Валидация дат
	if req.BirthYear != nil && req.DeathYear != nil && *req.BirthYear > *req.DeathYear {
		return httperr.New(fiber.StatusBadRequest, "Birth year cannot be after death year.")
	}

	if req.BirthYear != nil && *req.BirthYear > time.Now().Year() {
		return httperr.New(fiber.StatusBadRequest, "Birth year cannot be in the future.")
	}

	if req.DeathYear != nil && *req.DeathYear > time.Now().Year() {
		return httperr.New(fiber.StatusBadRequest, "Death year cannot be in the future.")
	}

	author, err := h.repo.CreateAuthor(c.Context(), postgres.CreateAuthorParams{
		FullName:  req.FullName,
		BirthYear: req.BirthYear,
		DeathYear: req.DeathYear,
		Biography: req.Biography,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Author with this name may already exist.")
		}
		log.Error().Err(err).Msg("Failed to create author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create author.", err.Error())
	}

	age := calculateAge(author.BirthYear, author.DeathYear)
	isAlive := author.DeathYear == nil
	booksCount := int64(0)

	response := AuthorResponse{
		ID:         author.ID,
		FullName:   author.FullName,
		BirthYear:  author.BirthYear,
		DeathYear:  author.DeathYear,
		Biography:  author.Biography,
		CreatedAt:  author.CreatedAt,
		BooksCount: &booksCount,
		Age:        age,
		IsAlive:    isAlive,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateAuthor обновляет информацию об авторе
func (h *Handler) updateAuthor(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can update authors.")
	}

	idStr := c.Params("id")
	authorID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format.")
	}

	var req UpdateAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию об авторе
	existingAuthor, err := h.repo.GetAuthorByID(c.Context(), authorID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Author not found.")
		}
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to get author for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve author.", err.Error())
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateAuthorParams{
		AuthorID: authorID,
	}

	if req.FullName != nil {
		if *req.FullName == "" {
			return httperr.New(fiber.StatusBadRequest, "Author full name cannot be empty.")
		}
		if len(*req.FullName) < 2 {
			return httperr.New(fiber.StatusBadRequest, "Author full name must be at least 2 characters long.")
		}
		updateParams.FullName = *req.FullName
	} else {
		updateParams.FullName = existingAuthor.FullName
	}

	if req.BirthYear != nil {
		updateParams.BirthYear = req.BirthYear
	} else {
		updateParams.BirthYear = existingAuthor.BirthYear
	}

	if req.DeathYear != nil {
		updateParams.DeathYear = req.DeathYear
	} else {
		updateParams.DeathYear = existingAuthor.DeathYear
	}

	if req.Biography != nil {
		updateParams.Biography = req.Biography
	} else {
		updateParams.Biography = existingAuthor.Biography
	}

	// Валидация дат
	if updateParams.BirthYear != nil && updateParams.DeathYear != nil && *updateParams.BirthYear > *updateParams.DeathYear {
		return httperr.New(fiber.StatusBadRequest, "Birth year cannot be after death year.")
	}

	if updateParams.BirthYear != nil && *updateParams.BirthYear > time.Now().Year() {
		return httperr.New(fiber.StatusBadRequest, "Birth year cannot be in the future.")
	}

	if updateParams.DeathYear != nil && *updateParams.DeathYear > time.Now().Year() {
		return httperr.New(fiber.StatusBadRequest, "Death year cannot be in the future.")
	}

	updatedAuthor, err := h.repo.UpdateAuthor(c.Context(), updateParams)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Author with this name may already exist.")
		}
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to update author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update author.", err.Error())
	}

	// Получаем количество книг автора
	books, err := h.repo.GetAuthorBooks(c.Context(), authorID)
	if err != nil {
		log.Warn().Err(err).Str("authorID", idStr).Msg("Failed to get author books count after update")
	}

	var booksCount *int64
	if books != nil {
		count := int64(len(books))
		booksCount = &count
	}

	age := calculateAge(updatedAuthor.BirthYear, updatedAuthor.DeathYear)
	isAlive := updatedAuthor.DeathYear == nil

	response := AuthorResponse{
		ID:         updatedAuthor.ID,
		FullName:   updatedAuthor.FullName,
		BirthYear:  updatedAuthor.BirthYear,
		DeathYear:  updatedAuthor.DeathYear,
		Biography:  updatedAuthor.Biography,
		CreatedAt:  updatedAuthor.CreatedAt,
		BooksCount: booksCount,
		Age:        age,
		IsAlive:    isAlive,
	}

	return c.JSON(response)
}

// deleteAuthor удаляет автора
func (h *Handler) deleteAuthor(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) {
		// library/server/internal/handler/authors.go
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators can delete authors.")
	}

	idStr := c.Params("id")
	authorID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format.")
	}

	// Проверяем, что автор существует
	_, err = h.repo.GetAuthorByID(c.Context(), authorID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Author not found.")
		}
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to get author for deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve author.", err.Error())
	}

	// Проверяем, есть ли у автора книги
	books, err := h.repo.GetAuthorBooks(c.Context(), authorID)
	if err != nil {
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to check author books before deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify author books.", err.Error())
	}

	if len(books) > 0 {
		return httperr.New(fiber.StatusConflict, "Cannot delete author who has books. Remove all book associations first.")
	}

	// Удаляем автора
	err = h.repo.DeleteAuthor(c.Context(), authorID)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusConflict, "Cannot delete author. Author has related records.")
		}
		log.Error().Err(err).Str("authorID", idStr).Msg("Failed to delete author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to delete author.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// getBookAuthors возвращает авторов книги
func (h *Handler) getBookAuthors(c *fiber.Ctx) error {
	bookIDStr := c.Params("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	// Проверяем, что книга существует
	_, err = h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to verify book existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify book.", err.Error())
	}

	authors, err := h.repo.GetBookAuthors(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to get book authors")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book authors.", err.Error())
	}

	response := make([]AuthorResponse, len(authors))
	for i, author := range authors {
		age := calculateAge(author.BirthYear, author.DeathYear)
		isAlive := author.DeathYear == nil

		response[i] = AuthorResponse{
			ID:        author.ID,
			FullName:  author.FullName,
			BirthYear: author.BirthYear,
			DeathYear: author.DeathYear,
			Biography: author.Biography,
			CreatedAt: author.CreatedAt,
			Age:       age,
			IsAlive:   isAlive,
		}
	}

	return c.JSON(fiber.Map{
		"authors": response,
		"book_id": bookID,
	})
}

// addBookAuthor добавляет автора к книге
func (h *Handler) addBookAuthor(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can manage book authors.")
	}

	bookIDStr := c.Params("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	var req AddBookAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Проверяем, что книга существует
	_, err = h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to verify book existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify book.", err.Error())
	}

	// Проверяем, что автор существует
	_, err = h.repo.GetAuthorByID(c.Context(), req.AuthorID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Author not found.")
		}
		log.Error().Err(err).Str("authorID", req.AuthorID.String()).Msg("Failed to verify author existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify author.", err.Error())
	}

	// Добавляем связь автора с книгой
	err = h.repo.AddBookAuthor(c.Context(), postgres.AddBookAuthorParams{
		BookID:   bookID,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid book ID or author ID.")
		}
		log.Error().Err(err).Str("bookID", bookIDStr).Str("authorID", req.AuthorID.String()).Msg("Failed to add book author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to add author to book.", err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Author successfully added to book",
		"book_id":   bookID,
		"author_id": req.AuthorID,
	})
}

// removeBookAuthor удаляет автора из книги
func (h *Handler) removeBookAuthor(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can manage book authors.")
	}

	bookIDStr := c.Params("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	authorIDStr := c.Params("authorId")
	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid author ID format.")
	}

	// Проверяем, что книга существует
	_, err = h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to verify book existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify book.", err.Error())
	}

	// Удаляем связь автора с книгой
	err = h.repo.RemoveBookAuthor(c.Context(), postgres.RemoveBookAuthorParams{
		BookID:   bookID,
		AuthorID: authorID,
	})
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIDStr).Str("authorID", authorIDStr).Msg("Failed to remove book author")
		return httperr.New(fiber.StatusInternalServerError, "Failed to remove author from book.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// removeAllBookAuthors удаляет всех авторов из книги
func (h *Handler) removeAllBookAuthors(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators can remove all authors from a book.")
	}

	bookIDStr := c.Params("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	// Проверяем, что книга существует
	_, err = h.repo.GetBookByID(c.Context(), bookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to verify book existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify book.", err.Error())
	}

	// Удаляем всех авторов книги
	err = h.repo.RemoveAllBookAuthors(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to remove all book authors")
		return httperr.New(fiber.StatusInternalServerError, "Failed to remove all authors from book.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

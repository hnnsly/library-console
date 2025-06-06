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

type IssueBookRequest struct {
	ReaderID string `json:"reader_id" validate:"required"`
	CopyCode string `json:"copy_code" validate:"required"`
	DueDays  int    `json:"due_days" validate:"required,min=1,max=365"`
}

type ReturnBookRequest struct {
	CopyCode string `json:"copy_code" validate:"required"`
}

func (h *Handler) getBooksToReturn(c *fiber.Ctx) error {
	books, err := h.repo.GetBooksToReturn(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get books to return")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve books to return")
	}

	return c.JSON(books)
}

func (h *Handler) getOverdueBooks(c *fiber.Ctx) error {
	books, err := h.repo.GetOverdueBooks(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get overdue books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve overdue books")
	}

	return c.JSON(books)
}

func (h *Handler) getRecentBookOperations(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "50")
	daysStr := c.Query("days", "7")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid limit parameter")
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid days parameter")
	}

	sinceDate := time.Now().AddDate(0, 0, -days)

	operations, err := h.repo.GetRecentBookOperations(c.Context(), postgres.GetRecentBookOperationsParams{
		LimitCount: int32(limit),
		SinceDate:  &sinceDate,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent book operations")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve recent book operations")
	}

	return c.JSON(operations)
}

func (h *Handler) issueBook(c *fiber.Ctx) error {
	var req IssueBookRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Parse reader ID
	readerID, err := uuid.Parse(req.ReaderID)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format")
	}

	// Get book copy by code
	bookCopy, err := h.repo.GetAvailableBookCopy(c.Context(), req.CopyCode)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Available book copy not found")
		}
		log.Error().Err(err).Str("copyCode", req.CopyCode).Msg("Failed to get available book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to check book availability")
	}

	// Check if reader has overdue books
	overdueCount, err := h.repo.CheckReaderOverdueBooks(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", req.ReaderID).Msg("Failed to check reader overdue books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to check reader status")
	}
	if overdueCount > 0 {
		return httperr.New(fiber.StatusForbidden, "Reader has overdue books")
	}

	// Get librarian ID from context
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return httperr.New(fiber.StatusUnauthorized, "User ID not found in context")
	}
	librarianID, err := uuid.Parse(userIDStr)
	if err != nil {
		return httperr.New(fiber.StatusInternalServerError, "Invalid user ID format")
	}

	dueDate := time.Now().AddDate(0, 0, req.DueDays)

	// Issue book
	issue, err := h.repo.IssueBook(c.Context(), postgres.IssueBookParams{
		ReaderID:    readerID,
		BookCopyID:  bookCopy.ID,
		DueDate:     dueDate,
		LibrarianID: &librarianID,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to issue book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to issue book")
	}

	// Update book copy status to issued
	err = h.repo.UpdateBookCopyStatus(c.Context(), postgres.UpdateBookCopyStatusParams{
		CopyID: bookCopy.ID,
		Status: postgres.NullBookStatus{
			BookStatus: postgres.BookStatusIssued,
			Valid:      true,
		},
	})
	if err != nil {
		log.Error().Err(err).Str("copyID", bookCopy.ID.String()).Msg("Failed to update book copy status")
		// Don't return error as the book was already issued
	}

	return c.Status(fiber.StatusCreated).JSON(issue)
}

func (h *Handler) returnBook(c *fiber.Ctx) error {
	var req ReturnBookRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Get book copy by code
	bookCopy, err := h.repo.GetBookCopyByCode(c.Context(), req.CopyCode)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found")
		}
		log.Error().Err(err).Str("copyCode", req.CopyCode).Msg("Failed to get book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to get book copy")
	}

	// Return book
	returnData, err := h.repo.ReturnBook(c.Context(), bookCopy.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "No active issue found for this book copy")
		}
		log.Error().Err(err).Str("copyCode", req.CopyCode).Msg("Failed to return book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to return book")
	}

	// Update book copy status to available
	err = h.repo.UpdateBookCopyStatus(c.Context(), postgres.UpdateBookCopyStatusParams{
		CopyID: bookCopy.ID,
		Status: postgres.NullBookStatus{
			BookStatus: postgres.BookStatusAvailable,
			Valid:      true,
		},
	})
	if err != nil {
		log.Error().Err(err).Str("copyID", bookCopy.ID.String()).Msg("Failed to update book copy status")
		// Don't return error as the book was already returned
	}

	return c.JSON(returnData)
}

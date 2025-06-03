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

type CreateLoanRequest struct {
	BookID      int `json:"book_id"`
	ReaderID    int `json:"reader_id"`
	LibrarianID int `json:"librarian_id"`
}

type CreateRenewalRequest struct {
	LoanHistoryID int       `json:"loan_history_id"`
	OldDueDate    time.Time `json:"old_due_date"`
	NewDueDate    time.Time `json:"new_due_date"`
	LibrarianID   int       `json:"librarian_id"`
	Reason        *string   `json:"reason"`
}

type ReturnBookRequest struct {
	LibrarianID *int `json:"librarian_id"`
}

type GetReaderLoanHistoryRequest struct {
	ReaderID   int   `json:"reader_id"`
	PageOffset int32 `json:"page_offset"`
	PageLimit  int32 `json:"page_limit"`
}

// checkLoanEligibility проверяет возможность выдачи книги читателю
func (h *Handler) checkLoanEligibility(c *fiber.Ctx) error {
	readerIDStr := c.Query("reader_id")
	bookIDStr := c.Query("book_id")

	if readerIDStr == "" || bookIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Reader ID and Book ID are required.")
	}

	readerID, err := strconv.Atoi(readerIDStr)
	if err != nil || readerID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID.")
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID.")
	}

	// TODO: Validate reader_id and book_id exist

	params := postgres.CheckLoanEligibilityParams{
		ReaderID: readerID,
		BookID:   bookID,
	}

	eligibility, err := h.repo.CheckLoanEligibility(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Int("readerID", readerID).Int("bookID", bookID).Msg("Failed to check loan eligibility")
		return httperr.New(fiber.StatusInternalServerError, "Failed to check loan eligibility.")
	}

	return c.JSON(eligibility)
}

// createLoan создает новую выдачу книги
func (h *Handler) createLoan(c *fiber.Ctx) error {
	req := new(CreateLoanRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: book_id, reader_id, librarian_id
	// TODO: Validate all IDs > 0 and exist
	// TODO: Check loan eligibility before creating

	params := postgres.CreateLoanParams{
		BookID:      req.BookID,
		ReaderID:    req.ReaderID,
		LibrarianID: req.LibrarianID,
	}

	loan, err := h.repo.CreateLoan(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create loan")
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid book, reader, or librarian ID.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create loan.")
	}

	return c.Status(fiber.StatusCreated).JSON(loan)
}

// getLoanByID получает информацию о выдаче по ID
func (h *Handler) getLoanByID(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	loan, err := h.repo.GetLoanByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Loan not found.")
		}
		log.Error().Err(err).Int64("loanID", id).Msg("Failed to get loan by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve loan.")
	}

	return c.JSON(loan)
}

// getReaderCurrentLoans получает текущие выдачи читателя
func (h *Handler) getReaderCurrentLoans(c *fiber.Ctx) error {
	readerIDStr := c.Params("reader_id")
	if readerIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Reader ID is required.")
	}

	readerID, err := strconv.Atoi(readerIDStr)
	if err != nil || readerID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID.")
	}

	// TODO: Validate reader_id exists

	loans, err := h.repo.GetReaderCurrentLoans(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Int("readerID", readerID).Msg("Failed to get reader current loans")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve current loans.")
	}

	if loans == nil {
		loans = []*postgres.GetReaderCurrentLoansRow{}
	}

	return c.JSON(loans)
}

// getReaderLoanHistory получает историю выдач читателя с пагинацией
func (h *Handler) getReaderLoanHistory(c *fiber.Ctx) error {
	req := new(GetReaderLoanHistoryRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate reader_id > 0 and exists
	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0

	if req.PageLimit == 0 {
		req.PageLimit = 20 // default limit
	}

	params := postgres.GetReaderLoanHistoryParams{
		ReaderID:   req.ReaderID,
		PageOffset: req.PageOffset,
		PageLimit:  req.PageLimit,
	}

	history, err := h.repo.GetReaderLoanHistory(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Int("readerID", req.ReaderID).Msg("Failed to get reader loan history")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve loan history.")
	}

	if history == nil {
		history = []*postgres.GetReaderLoanHistoryRow{}
	}

	return c.JSON(history)
}

// getOverdueBooks получает список просроченных книг
func (h *Handler) getOverdueBooks(c *fiber.Ctx) error {
	limit := int32(50) // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// TODO: Validate limit > 0 and <= 200

	overdueBooks, err := h.repo.GetOverdueBooks(c.Context(), limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get overdue books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve overdue books.")
	}

	if overdueBooks == nil {
		overdueBooks = []*postgres.GetOverdueBooksRow{}
	}

	return c.JSON(overdueBooks)
}

// getBooksDueToday получает список книг, которые должны быть возвращены сегодня
func (h *Handler) getBooksDueToday(c *fiber.Ctx) error {
	booksDue, err := h.repo.GetBooksDueToday(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get books due today")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve books due today.")
	}

	if booksDue == nil {
		booksDue = []*postgres.GetBooksDueTodayRow{}
	}

	return c.JSON(booksDue)
}

// getActiveLoansByBook получает активные выдачи конкретной книги
func (h *Handler) getActiveLoansByBook(c *fiber.Ctx) error {
	bookIDStr := c.Params("book_id")
	if bookIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Book ID is required.")
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil || bookID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID.")
	}

	// TODO: Validate book_id exists

	loans, err := h.repo.GetActiveLoansByBook(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Int("bookID", bookID).Msg("Failed to get active loans by book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve active loans.")
	}

	if loans == nil {
		loans = []*postgres.GetActiveLoansByBookRow{}
	}

	return c.JSON(loans)
}

// returnBook возвращает книгу
func (h *Handler) returnBook(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	req := new(ReturnBookRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate loan exists and is active
	// TODO: Validate librarian_id if provided

	params := postgres.ReturnBookParams{
		LoanID:      id,
		LibrarianID: req.LibrarianID,
	}

	err = h.repo.ReturnBook(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Active loan not found.")
		}
		log.Error().Err(err).Int64("loanID", id).Msg("Failed to return book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to return book.")
	}

	return c.JSON(fiber.Map{"message": "Book successfully returned"})
}

// renewLoan продлевает срок выдачи книги
func (h *Handler) renewLoan(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.RenewLoan(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Active loan not found.")
		}
		log.Error().Err(err).Int64("loanID", id).Msg("Failed to renew loan")
		return httperr.New(fiber.StatusInternalServerError, "Failed to renew loan.")
	}

	return c.JSON(fiber.Map{"message": "Loan successfully renewed"})
}

// createRenewal создает запись о продлении
func (h *Handler) createRenewal(c *fiber.Ctx) error {
	req := new(CreateRenewalRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: loan_history_id, old_due_date, new_due_date, librarian_id
	// TODO: Validate loan_history_id exists and is active
	// TODO: Validate new_due_date > old_due_date
	// TODO: Validate librarian_id exists

	params := postgres.CreateRenewalParams{
		LoanHistoryID: req.LoanHistoryID,
		OldDueDate:    req.OldDueDate,
		NewDueDate:    req.NewDueDate,
		LibrarianID:   req.LibrarianID,
		Reason:        req.Reason,
	}

	err := h.repo.CreateRenewal(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create renewal")
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid loan history or librarian ID.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create renewal.")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Renewal successfully created"})
}

// markLoanAsLost помечает выдачу как потерянную
func (h *Handler) markLoanAsLost(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.MarkLoanAsLost(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Loan not found.")
		}
		log.Error().Err(err).Int64("loanID", id).Msg("Failed to mark loan as lost")
		return httperr.New(fiber.StatusInternalServerError, "Failed to mark loan as lost.")
	}

	return c.JSON(fiber.Map{"message": "Loan successfully marked as lost"})
}

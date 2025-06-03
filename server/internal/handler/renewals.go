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

type CreateRenewalRecordRequest struct {
	LoanHistoryID int       `json:"loan_history_id"`
	OldDueDate    time.Time `json:"old_due_date"`
	NewDueDate    time.Time `json:"new_due_date"`
	LibrarianID   int       `json:"librarian_id"`
	Reason        *string   `json:"reason"`
}

type GetRenewalsByDateRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// createRenewalRecord создает запись о продлении
func (h *Handler) createRenewalRecord(c *fiber.Ctx) error {
	req := new(CreateRenewalRecordRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: loan_history_id, old_due_date, new_due_date, librarian_id
	// TODO: Validate loan_history_id exists and is active
	// TODO: Validate new_due_date > old_due_date
	// TODO: Validate new_due_date is not too far in future
	// TODO: Validate librarian_id exists
	// TODO: Validate reason max length 500 if provided

	params := postgres.CreateRenewalRecordParams{
		LoanHistoryID: req.LoanHistoryID,
		OldDueDate:    req.OldDueDate,
		NewDueDate:    req.NewDueDate,
		LibrarianID:   req.LibrarianID,
		Reason:        req.Reason,
	}

	renewal, err := h.repo.CreateRenewalRecord(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create renewal record")
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid loan history or librarian ID.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create renewal record.")
	}

	return c.Status(fiber.StatusCreated).JSON(renewal)
}

// getRenewalsForLoan получает все продления для конкретной выдачи
func (h *Handler) getRenewalsForLoan(c *fiber.Ctx) error {
	loanHistoryIDStr := c.Params("loan_id")
	if loanHistoryIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Loan history ID is required.")
	}

	loanHistoryID, err := strconv.Atoi(loanHistoryIDStr)
	if err != nil || loanHistoryID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid loan history ID.")
	}

	// TODO: Validate loan_history_id exists

	renewals, err := h.repo.GetRenewalsForLoan(c.Context(), loanHistoryID)
	if err != nil {
		log.Error().Err(err).Int("loanHistoryID", loanHistoryID).Msg("Failed to get renewals for loan")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve renewals.")
	}

	if renewals == nil {
		renewals = []*postgres.GetRenewalsForLoanRow{}
	}

	return c.JSON(renewals)
}

// getRenewalsByDate получает продления за период
func (h *Handler) getRenewalsByDate(c *fiber.Ctx) error {
	req := new(GetRenewalsByDateRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: start_date, end_date
	// TODO: Validate end_date >= start_date
	// TODO: Validate date range is not too large (max 1 year)

	params := postgres.GetRenewalsByDateParams{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	renewals, err := h.repo.GetRenewalsByDate(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get renewals by date")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve renewals.")
	}

	if renewals == nil {
		renewals = []*postgres.GetRenewalsByDateRow{}
	}

	return c.JSON(renewals)
}

// getMostRenewedBooks получает список самых продлеваемых книг
func (h *Handler) getMostRenewedBooks(c *fiber.Ctx) error {
	limit := int32(20) // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// TODO: Validate limit > 0 and <= 100

	books, err := h.repo.GetMostRenewedBooks(c.Context(), limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get most renewed books")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve most renewed books.")
	}

	if books == nil {
		books = []*postgres.GetMostRenewedBooksRow{}
	}

	return c.JSON(books)
}

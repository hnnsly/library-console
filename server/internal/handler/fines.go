package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/govalues/decimal"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type CreateFineRequest struct {
	LoanHistoryID int             `json:"loan_history_id"`
	ReaderID      int             `json:"reader_id"`
	FineType      string          `json:"fine_type"`
	Amount        decimal.Decimal `json:"amount"`
	Description   *string         `json:"description"`
	LibrarianID   int             `json:"librarian_id"`
}

type CalculateOverdueFineRequest struct {
	LoanHistoryID int64           `json:"loan_history_id"`
	DailyFineRate decimal.Decimal `json:"daily_fine_rate"`
}

// createFine создает новый штраф
func (h *Handler) createFine(c *fiber.Ctx) error {
	req := new(CreateFineRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: loan_history_id, reader_id, fine_type, amount, librarian_id
	// TODO: Validate loan_history_id > 0 and exists
	// TODO: Validate reader_id > 0 and exists
	// TODO: Validate librarian_id > 0 and exists
	// TODO: Validate fine_type is one of: overdue, lost, damage, other
	// TODO: Validate amount > 0
	// TODO: Validate description max length 500 if provided

	params := postgres.CreateFineParams{
		LoanHistoryID: req.LoanHistoryID,
		ReaderID:      req.ReaderID,
		FineType:      req.FineType,
		Amount:        req.Amount,
		Description:   req.Description,
		LibrarianID:   req.LibrarianID,
	}

	fine, err := h.repo.CreateFine(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create fine")
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid loan history, reader, or librarian ID.")
		}
		return httperr.New(fiber.StatusInternalServerError, "Failed to create fine.")
	}

	return c.Status(fiber.StatusCreated).JSON(fine)
}

// calculateOverdueFine рассчитывает штраф за просрочку
func (h *Handler) calculateOverdueFine(c *fiber.Ctx) error {
	req := new(CalculateOverdueFineRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate loan_history_id > 0 and exists
	// TODO: Validate daily_fine_rate > 0

	// Преобразуем decimal.Decimal в pgtype.Numeric
	var dailyFineRate pgtype.Numeric
	if err := dailyFineRate.Scan(req.DailyFineRate.String()); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid daily fine rate format.")
	}

	params := postgres.CalculateOverdueFineParams{
		DailyFineRate: dailyFineRate,
		LoanHistoryID: req.LoanHistoryID,
	}

	result, err := h.repo.CalculateOverdueFine(c.Context(), params)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Loan history not found.")
		}
		log.Error().Err(err).Int64("loanHistoryID", req.LoanHistoryID).Msg("Failed to calculate overdue fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to calculate overdue fine.")
	}

	// Преобразуем pgtype.Numeric в decimal.Decimal для ответа
	var calculatedFine decimal.Decimal
	if result.CalculatedFine.Valid {
		calculatedFine, _ = decimal.Parse(result.CalculatedFine.Int.String())
	}

	response := map[string]interface{}{
		"loan_history_id": result.LoanHistoryID,
		"overdue_days":    result.OverdueDays,
		"calculated_fine": calculatedFine,
	}

	return c.JSON(response)
}

// payFine помечает штраф как оплаченный
func (h *Handler) payFine(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.PayFine(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Fine not found.")
		}
		log.Error().Err(err).Int64("fineID", id).Msg("Failed to pay fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to pay fine.")
	}

	return c.JSON(fiber.Map{"message": "Fine successfully paid"})
}

// waiveFine аннулирует штраф
func (h *Handler) waiveFine(c *fiber.Ctx) error {
	id, err := parseID(c, "id")
	if err != nil {
		return err
	}

	err = h.repo.WaiveFine(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") {
			return httperr.New(fiber.StatusNotFound, "Fine not found.")
		}
		log.Error().Err(err).Int64("fineID", id).Msg("Failed to waive fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to waive fine.")
	}

	return c.JSON(fiber.Map{"message": "Fine successfully waived"})
}

// getReaderFines получает все штрафы читателя
func (h *Handler) getReaderFines(c *fiber.Ctx) error {
	readerIDStr := c.Params("reader_id")
	if readerIDStr == "" {
		return httperr.New(fiber.StatusBadRequest, "Reader ID is required.")
	}

	readerID, err := strconv.Atoi(readerIDStr)
	if err != nil || readerID <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID.")
	}

	// TODO: Validate reader_id exists

	fines, err := h.repo.GetReaderFines(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Int("readerID", readerID).Msg("Failed to get reader fines")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader fines.")
	}

	if fines == nil {
		fines = []*postgres.GetReaderFinesRow{}
	}

	return c.JSON(fines)
}

// getUnpaidFines получает все неоплаченные штрафы
func (h *Handler) getUnpaidFines(c *fiber.Ctx) error {
	fines, err := h.repo.GetUnpaidFines(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get unpaid fines")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve unpaid fines.")
	}

	if fines == nil {
		fines = []*postgres.GetUnpaidFinesRow{}
	}

	return c.JSON(fines)
}

// getDebtorReaders получает список читателей-должников
func (h *Handler) getDebtorReaders(c *fiber.Ctx) error {
	debtors, err := h.repo.GetDebtorReaders(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get debtor readers")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve debtor readers.")
	}

	if debtors == nil {
		debtors = []*postgres.GetDebtorReadersRow{}
	}

	return c.JSON(debtors)
}

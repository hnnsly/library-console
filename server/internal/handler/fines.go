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

// Request/Response structs для fines
type CreateFineRequest struct {
	ReaderID    uuid.UUID       `json:"reader_id" validate:"required"`
	BookIssueID *uuid.UUID      `json:"book_issue_id"`
	Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Reason      string          `json:"reason" validate:"required,min=3,max=500"`
	FineDate    *time.Time      `json:"fine_date"`
	LibrarianID *uuid.UUID      `json:"librarian_id"`
}

type UpdateFineRequest struct {
	Amount     *decimal.Decimal `json:"amount" validate:"omitempty,gt=0"`
	Reason     *string          `json:"reason" validate:"omitempty,min=3,max=500"`
	PaidDate   *time.Time       `json:"paid_date"`
	PaidAmount *decimal.Decimal `json:"paid_amount" validate:"omitempty,gte=0"`
	IsPaid     *bool            `json:"is_paid"`
}

type PayFineRequest struct {
	PaidAmount decimal.Decimal `json:"paid_amount" validate:"required,gt=0"`
	PaidDate   *time.Time      `json:"paid_date"`
}

type FineResponse struct {
	ID          uuid.UUID       `json:"id"`
	ReaderID    uuid.UUID       `json:"reader_id"`
	BookIssueID *uuid.UUID      `json:"book_issue_id"`
	Amount      decimal.Decimal `json:"amount"`
	Reason      string          `json:"reason"`
	FineDate    *time.Time      `json:"fine_date"`
	PaidDate    *time.Time      `json:"paid_date"`
	PaidAmount  decimal.Decimal `json:"paid_amount"`
	IsPaid      *bool           `json:"is_paid"`
	LibrarianID *uuid.UUID      `json:"librarian_id"`
	CreatedAt   *time.Time      `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
	// Extended fields from joins
	ReaderName    *string    `json:"reader_name,omitempty"`
	TicketNumber  *string    `json:"ticket_number,omitempty"`
	LibrarianName *string    `json:"librarian_name,omitempty"`
	BookTitle     *string    `json:"book_title,omitempty"`
	IssueDate     *time.Time `json:"issue_date,omitempty"`
	DueDate       *time.Time `json:"due_date,omitempty"`
}

type FineStatisticsResponse struct {
	TotalFines  int64           `json:"total_fines"`
	PaidFines   int64           `json:"paid_fines"`
	UnpaidFines int64           `json:"unpaid_fines"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	TotalPaid   decimal.Decimal `json:"total_paid"`
	TotalDebt   decimal.Decimal `json:"total_debt"`
}

type ReaderDebtResponse struct {
	ReaderID    uuid.UUID       `json:"reader_id"`
	TotalDebt   decimal.Decimal `json:"total_debt"`
	UnpaidFines int64           `json:"unpaid_fines"`
}

// listFines возвращает список всех штрафов (только для администраторов/библиотекарей)
func (h *Handler) listFines(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view all fines.")
	}

	// Для простоты возвращаем неоплаченные штрафы
	return h.getUnpaidFines(c)
}

// getUnpaidFines возвращает список всех неоплаченных штрафов
func (h *Handler) getUnpaidFines(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view unpaid fines.")
	}

	fines, err := h.repo.GetAllUnpaidFines(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get unpaid fines")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve unpaid fines.", err.Error())
	}

	response := make([]FineResponse, len(fines))
	for i, fine := range fines {
		response[i] = FineResponse{
			ID:           fine.ID,
			ReaderID:     fine.ReaderID,
			BookIssueID:  fine.BookIssueID,
			Amount:       fine.Amount,
			Reason:       fine.Reason,
			FineDate:     fine.FineDate,
			PaidDate:     fine.PaidDate,
			PaidAmount:   fine.PaidAmount,
			IsPaid:       fine.IsPaid,
			LibrarianID:  fine.LibrarianID,
			CreatedAt:    fine.CreatedAt,
			UpdatedAt:    fine.UpdatedAt,
			ReaderName:   &fine.ReaderName,
			TicketNumber: &fine.TicketNumber,
			BookTitle:    fine.BookTitle,
		}
	}

	return c.JSON(fiber.Map{"fines": response})
}

// getFineStatistics возвращает статистику по штрафам за указанный период
func (h *Handler) getFineStatistics(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view statistics.")
	}

	// Параметры периода (по умолчанию - последние 30 дней)
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	var fromDate, toDate *time.Time
	now := time.Now()

	if fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = &parsed
		} else {
			return httperr.New(fiber.StatusBadRequest, "Invalid from_date format. Use YYYY-MM-DD.")
		}
	} else {
		defaultFrom := now.AddDate(0, 0, -30)
		fromDate = &defaultFrom
	}

	if toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = &parsed
		} else {
			return httperr.New(fiber.StatusBadRequest, "Invalid to_date format. Use YYYY-MM-DD.")
		}
	} else {
		toDate = &now
	}

	statistics, err := h.repo.GetFineStatistics(c.Context(), postgres.GetFineStatisticsParams{
		FromDate: fromDate,
		ToDate:   toDate,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get fine statistics")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve fine statistics.", err.Error())
	}

	// Конвертируем interface{} в decimal.Decimal
	totalAmount, _ := decimal.Parse(statistics.TotalAmount.(string))
	totalPaid, _ := decimal.Parse(statistics.TotalPaid.(string))
	totalDebt, _ := decimal.Parse(statistics.TotalDebt.(string))

	response := FineStatisticsResponse{
		TotalFines:  statistics.TotalFines,
		PaidFines:   statistics.PaidFines,
		UnpaidFines: statistics.UnpaidFines,
		TotalAmount: totalAmount,
		TotalPaid:   totalPaid,
		TotalDebt:   totalDebt,
	}

	return c.JSON(response)
}

// getMyFines возвращает штрафы текущего читателя
func (h *Handler) getMyFines(c *fiber.Ctx) error {
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
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get reader for fines")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	return h.getFinesByReaderInternal(c, reader.ID)
}

// getFinesByReader возвращает штрафы читателя по ID
func (h *Handler) getFinesByReader(c *fiber.Ctx) error {
	readerIDStr := c.Params("readerId")
	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format.")
	}

	// Проверка доступа: читатель может видеть только свои штрафы
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != readerID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own fines.")
		}
	}

	return h.getFinesByReaderInternal(c, readerID)
}

// getFinesByReaderInternal - внутренний метод для получения штрафов читателя
func (h *Handler) getFinesByReaderInternal(c *fiber.Ctx, readerID uuid.UUID) error {
	fines, err := h.repo.GetFinesByReader(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", readerID.String()).Msg("Failed to get fines by reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve fines.", err.Error())
	}

	response := make([]FineResponse, len(fines))
	for i, fine := range fines {
		response[i] = FineResponse{
			ID:          fine.ID,
			ReaderID:    fine.ReaderID,
			BookIssueID: fine.BookIssueID,
			Amount:      fine.Amount,
			Reason:      fine.Reason,
			FineDate:    fine.FineDate,
			PaidDate:    fine.PaidDate,
			PaidAmount:  fine.PaidAmount,
			IsPaid:      fine.IsPaid,
			LibrarianID: fine.LibrarianID,
			CreatedAt:   fine.CreatedAt,
			UpdatedAt:   fine.UpdatedAt,
			BookTitle:   fine.BookTitle,
			IssueDate:   fine.IssueDate,
			DueDate:     fine.DueDate,
		}
	}

	return c.JSON(fiber.Map{"fines": response})
}

// getUnpaidFinesByReader возвращает неоплаченные штрафы читателя
func (h *Handler) getUnpaidFinesByReader(c *fiber.Ctx) error {
	readerIDStr := c.Params("readerId")
	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format.")
	}

	// Проверка доступа
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != readerID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own fines.")
		}
	}

	fines, err := h.repo.GetUnpaidFinesByReader(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", readerIDStr).Msg("Failed to get unpaid fines by reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve unpaid fines.", err.Error())
	}

	response := make([]FineResponse, len(fines))
	for i, fine := range fines {
		response[i] = FineResponse{
			ID:          fine.ID,
			ReaderID:    fine.ReaderID,
			BookIssueID: fine.BookIssueID,
			Amount:      fine.Amount,
			Reason:      fine.Reason,
			FineDate:    fine.FineDate,
			PaidDate:    fine.PaidDate,
			PaidAmount:  fine.PaidAmount,
			IsPaid:      fine.IsPaid,
			LibrarianID: fine.LibrarianID,
			CreatedAt:   fine.CreatedAt,
			UpdatedAt:   fine.UpdatedAt,
			BookTitle:   fine.BookTitle,
			IssueDate:   fine.IssueDate,
			DueDate:     fine.DueDate,
		}
	}

	return c.JSON(fiber.Map{"unpaid_fines": response})
}

// getReaderDebt возвращает общий долг читателя
func (h *Handler) getReaderDebt(c *fiber.Ctx) error {
	readerIDStr := c.Params("readerId")
	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid reader ID format.")
	}

	// Проверка доступа
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != readerID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own debt.")
		}
	}

	totalDebtInterface, err := h.repo.GetTotalDebtByReader(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", readerIDStr).Msg("Failed to get reader debt")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader debt.", err.Error())
	}

	// Получаем количество неоплаченных штрафов
	unpaidFines, err := h.repo.GetUnpaidFinesByReader(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", readerIDStr).Msg("Failed to get unpaid fines count")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve unpaid fines count.", err.Error())
	}

	totalDebt, _ := decimal.Parse(totalDebtInterface.(string))

	response := ReaderDebtResponse{
		ReaderID:    readerID,
		TotalDebt:   totalDebt,
		UnpaidFines: int64(len(unpaidFines)),
	}

	return c.JSON(response)
}

// getFineByID возвращает штраф по ID
func (h *Handler) getFineByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	fineID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid fine ID format.")
	}

	fine, err := h.repo.GetFineByID(c.Context(), fineID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Fine not found.")
		}
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to get fine by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve fine.", err.Error())
	}

	// Проверка доступа для читателей
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != fine.ReaderID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own fines.")
		}
	}

	response := FineResponse{
		ID:            fine.ID,
		ReaderID:      fine.ReaderID,
		BookIssueID:   fine.BookIssueID,
		Amount:        fine.Amount,
		Reason:        fine.Reason,
		FineDate:      fine.FineDate,
		PaidDate:      fine.PaidDate,
		PaidAmount:    fine.PaidAmount,
		IsPaid:        fine.IsPaid,
		LibrarianID:   fine.LibrarianID,
		CreatedAt:     fine.CreatedAt,
		UpdatedAt:     fine.UpdatedAt,
		ReaderName:    &fine.ReaderName,
		TicketNumber:  &fine.TicketNumber,
		LibrarianName: fine.LibrarianName,
	}

	return c.JSON(response)
}

// createFine создает новый штраф
func (h *Handler) createFine(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can create fines.")
	}

	var req CreateFineRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация
	if req.Amount.Sign() <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Fine amount must be greater than 0.")
	}
	if req.Reason == "" {
		return httperr.New(fiber.StatusBadRequest, "Fine reason is required.")
	}

	// Если дата штрафа не указана, используем текущую дату
	if req.FineDate == nil {
		now := time.Now()
		req.FineDate = &now
	}

	// Если библиотекарь не указан, берем текущего пользователя
	if req.LibrarianID == nil {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)
		req.LibrarianID = &userID
	}

	fine, err := h.repo.CreateFine(c.Context(), postgres.CreateFineParams{
		ReaderID:    req.ReaderID,
		BookIssueID: req.BookIssueID,
		Amount:      req.Amount,
		Reason:      req.Reason,
		FineDate:    req.FineDate,
		LibrarianID: req.LibrarianID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid reader ID or book issue ID.")
		}
		log.Error().Err(err).Msg("Failed to create fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create fine.", err.Error())
	}

	response := FineResponse{
		ID:          fine.ID,
		ReaderID:    fine.ReaderID,
		BookIssueID: fine.BookIssueID,
		Amount:      fine.Amount,
		Reason:      fine.Reason,
		FineDate:    fine.FineDate,
		PaidDate:    fine.PaidDate,
		PaidAmount:  fine.PaidAmount,
		IsPaid:      fine.IsPaid,
		LibrarianID: fine.LibrarianID,
		CreatedAt:   fine.CreatedAt,
		UpdatedAt:   fine.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateFine обновляет штраф
func (h *Handler) updateFine(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can update fines.")
	}

	idStr := c.Params("id")
	fineID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid fine ID format.")
	}

	var req UpdateFineRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию о штрафе
	existingFine, err := h.repo.GetFineByID(c.Context(), fineID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Fine not found.")
		}
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to get fine for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve fine.", err.Error())
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateFineParams{
		FineID: fineID,
	}

	if req.Amount != nil {
		if req.Amount.Sign() <= 0 {
			return httperr.New(fiber.StatusBadRequest, "Fine amount must be greater than 0.")
		}
		updateParams.Amount = *req.Amount
	} else {
		updateParams.Amount = existingFine.Amount
	}

	if req.Reason != nil {
		if *req.Reason == "" {
			return httperr.New(fiber.StatusBadRequest, "Fine reason cannot be empty.")
		}
		updateParams.Reason = *req.Reason
	} else {
		updateParams.Reason = existingFine.Reason
	}

	if req.PaidDate != nil {
		updateParams.PaidDate = req.PaidDate
	} else {
		updateParams.PaidDate = existingFine.PaidDate
	}

	if req.PaidAmount != nil {
		if req.PaidAmount.Sign() < 0 {
			return httperr.New(fiber.StatusBadRequest, "Paid amount cannot be negative.")
		}
		updateParams.PaidAmount = *req.PaidAmount
	} else {
		updateParams.PaidAmount = existingFine.PaidAmount
	}

	if req.IsPaid != nil {
		updateParams.IsPaid = req.IsPaid
	} else {
		updateParams.IsPaid = existingFine.IsPaid
	}

	updatedFine, err := h.repo.UpdateFine(c.Context(), updateParams)
	if err != nil {
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to update fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update fine.", err.Error())
	}

	response := FineResponse{
		ID:          updatedFine.ID,
		ReaderID:    updatedFine.ReaderID,
		BookIssueID: updatedFine.BookIssueID,
		Amount:      updatedFine.Amount,
		Reason:      updatedFine.Reason,
		FineDate:    updatedFine.FineDate,
		PaidDate:    updatedFine.PaidDate,
		PaidAmount:  updatedFine.PaidAmount,
		IsPaid:      updatedFine.IsPaid,
		LibrarianID: updatedFine.LibrarianID,
		CreatedAt:   updatedFine.CreatedAt,
		UpdatedAt:   updatedFine.UpdatedAt,
	}

	return c.JSON(response)
}

// payFine обрабатывает оплату штрафа
func (h *Handler) payFine(c *fiber.Ctx) error {
	idStr := c.Params("id")
	fineID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid fine ID format.")
	}

	var req PayFineRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация суммы оплаты
	if req.PaidAmount.Sign() <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Paid amount must be greater than 0.")
	}

	// Проверяем, что штраф существует
	existingFine, err := h.repo.GetFineByID(c.Context(), fineID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Fine not found.")
		}
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to get fine for payment")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve fine.", err.Error())
	}

	// Проверка доступа для читателей
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != existingFine.ReaderID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only pay your own fines.")
		}
	}

	// Если дата оплаты не указана, используем текущую дату
	if req.PaidDate == nil {
		now := time.Now()
		req.PaidDate = &now
	}

	// Обрабатываем оплату
	err = h.repo.PayFine(c.Context(), postgres.PayFineParams{
		FineID:     fineID,
		PaidDate:   req.PaidDate,
		PaidAmount: req.PaidAmount,
	})
	if err != nil {
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to pay fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to process fine payment.", err.Error())
	}

	// Получаем обновленную информацию о штрафе
	updatedFine, err := h.repo.GetFineByID(c.Context(), fineID)
	if err != nil {
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to get updated fine")
		return httperr.New(fiber.StatusInternalServerError, "Payment processed but failed to retrieve updated fine.", err.Error())
	}

	response := FineResponse{
		ID:            updatedFine.ID,
		ReaderID:      updatedFine.ReaderID,
		BookIssueID:   updatedFine.BookIssueID,
		Amount:        updatedFine.Amount,
		Reason:        updatedFine.Reason,
		FineDate:      updatedFine.FineDate,
		PaidDate:      updatedFine.PaidDate,
		PaidAmount:    updatedFine.PaidAmount,
		IsPaid:        updatedFine.IsPaid,
		LibrarianID:   updatedFine.LibrarianID,
		CreatedAt:     updatedFine.CreatedAt,
		UpdatedAt:     updatedFine.UpdatedAt,
		ReaderName:    &updatedFine.ReaderName,
		TicketNumber:  &updatedFine.TicketNumber,
		LibrarianName: updatedFine.LibrarianName,
	}

	return c.JSON(response)
}

// deleteFine удаляет штраф
func (h *Handler) deleteFine(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators can delete fines.")
	}

	idStr := c.Params("id")
	fineID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid fine ID format.")
	}

	// Проверяем, что штраф существует
	_, err = h.repo.GetFineByID(c.Context(), fineID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Fine not found.")
		}
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to get fine for deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve fine.", err.Error())
	}

	// Удаляем штраф
	err = h.repo.DeleteFine(c.Context(), fineID)
	if err != nil {
		log.Error().Err(err).Str("fineID", idStr).Msg("Failed to delete fine")
		return httperr.New(fiber.StatusInternalServerError, "Failed to delete fine.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

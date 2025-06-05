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

// Request/Response structs для issues
type CreateIssueRequest struct {
	ReaderID    uuid.UUID  `json:"reader_id" validate:"required"`
	BookCopyID  uuid.UUID  `json:"book_copy_id" validate:"required"`
	IssueDate   *time.Time `json:"issue_date"`
	DueDate     time.Time  `json:"due_date" validate:"required"`
	LibrarianID *uuid.UUID `json:"librarian_id"`
	Notes       *string    `json:"notes" validate:"omitempty,max=1000"`
}

type UpdateIssueRequest struct {
	DueDate       *time.Time `json:"due_date"`
	ReturnDate    *time.Time `json:"return_date"`
	ExtendedCount *int       `json:"extended_count" validate:"omitempty,min=0,max=10"`
	Notes         *string    `json:"notes" validate:"omitempty,max=1000"`
}

type ExtendIssueRequest struct {
	ExtensionDays int     `json:"extension_days" validate:"required,min=1,max=30"`
	Notes         *string `json:"notes" validate:"omitempty,max=500"`
}

type ReturnBookRequest struct {
	ReturnDate *time.Time `json:"return_date"`
	Notes      *string    `json:"notes" validate:"omitempty,max=500"`
}

type IssueResponse struct {
	ID            uuid.UUID  `json:"id"`
	ReaderID      uuid.UUID  `json:"reader_id"`
	BookCopyID    uuid.UUID  `json:"book_copy_id"`
	IssueDate     *time.Time `json:"issue_date"`
	DueDate       time.Time  `json:"due_date"`
	ReturnDate    *time.Time `json:"return_date"`
	ExtendedCount *int       `json:"extended_count"`
	LibrarianID   *uuid.UUID `json:"librarian_id"`
	Notes         *string    `json:"notes"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	// Extended fields
	ReaderName    *string `json:"reader_name,omitempty"`
	TicketNumber  *string `json:"ticket_number,omitempty"`
	BookTitle     *string `json:"book_title,omitempty"`
	CopyCode      *string `json:"copy_code,omitempty"`
	LibrarianName *string `json:"librarian_name,omitempty"`
	OverdueDays   *int32  `json:"overdue_days,omitempty"`
	IsOverdue     bool    `json:"is_overdue"`
	IsActive      bool    `json:"is_active"`
}

type IssuesListResponse struct {
	Issues []IssueResponse `json:"issues"`
	Total  int64           `json:"total"`
	Limit  int32           `json:"limit"`
	Offset int32           `json:"offset"`
}

type IssueStatisticsResponse struct {
	TotalActiveIssues int64 `json:"total_active_issues"`
	OverdueIssues     int64 `json:"overdue_issues"`
	DueSoonIssues     int64 `json:"due_soon_issues"`
}

// listIssues возвращает список всех выдач (только для администраторов/библиотекарей)
func (h *Handler) listIssues(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view all issues.")
	}

	// Для простоты возвращаем активные выдачи
	return h.getActiveIssues(c)
}

// getActiveIssues возвращает все активные выдачи
func (h *Handler) getActiveIssues(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view all active issues.")
	}

	issues, err := h.repo.GetAllActiveIssues(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all active issues")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve active issues.", err.Error())
	}

	response := make([]IssueResponse, len(issues))
	for i, issue := range issues {
		isOverdue := issue.OverdueDays > 0

		response[i] = IssueResponse{
			ID:            issue.ID,
			ReaderID:      issue.ReaderID,
			BookCopyID:    issue.BookCopyID,
			IssueDate:     issue.IssueDate,
			DueDate:       issue.DueDate,
			ReturnDate:    issue.ReturnDate,
			ExtendedCount: issue.ExtendedCount,
			LibrarianID:   issue.LibrarianID,
			Notes:         issue.Notes,
			CreatedAt:     issue.CreatedAt,
			UpdatedAt:     issue.UpdatedAt,
			ReaderName:    &issue.ReaderName,
			TicketNumber:  &issue.TicketNumber,
			BookTitle:     &issue.BookTitle,
			CopyCode:      &issue.CopyCode,
			OverdueDays:   &issue.OverdueDays,
			IsOverdue:     isOverdue,
			IsActive:      true,
		}
	}

	return c.JSON(fiber.Map{"issues": response})
}

// getOverdueIssues возвращает просроченные выдачи
func (h *Handler) getOverdueIssues(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view overdue issues.")
	}

	issues, err := h.repo.GetOverdueIssues(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get overdue issues")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve overdue issues.", err.Error())
	}

	response := make([]IssueResponse, len(issues))
	for i, issue := range issues {
		response[i] = IssueResponse{
			ID:            issue.ID,
			ReaderID:      issue.ReaderID,
			BookCopyID:    issue.BookCopyID,
			IssueDate:     issue.IssueDate,
			DueDate:       issue.DueDate,
			ReturnDate:    issue.ReturnDate,
			ExtendedCount: issue.ExtendedCount,
			LibrarianID:   issue.LibrarianID,
			Notes:         issue.Notes,
			CreatedAt:     issue.CreatedAt,
			UpdatedAt:     issue.UpdatedAt,
			ReaderName:    &issue.ReaderName,
			TicketNumber:  &issue.TicketNumber,
			BookTitle:     &issue.BookTitle,
			CopyCode:      &issue.CopyCode,
			OverdueDays:   &issue.OverdueDays,
			IsOverdue:     true,
			IsActive:      true,
		}
	}

	return c.JSON(fiber.Map{"overdue_issues": response})
}

// getIssuesDueSoon возвращает выдачи, которые скоро истекают
func (h *Handler) getIssuesDueSoon(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view issues due soon.")
	}

	issues, err := h.repo.GetIssueDueSoon(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get issues due soon")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve issues due soon.", err.Error())
	}

	response := make([]IssueResponse, len(issues))
	for i, issue := range issues {
		response[i] = IssueResponse{
			ID:            issue.ID,
			ReaderID:      issue.ReaderID,
			BookCopyID:    issue.BookCopyID,
			IssueDate:     issue.IssueDate,
			DueDate:       issue.DueDate,
			ReturnDate:    issue.ReturnDate,
			ExtendedCount: issue.ExtendedCount,
			LibrarianID:   issue.LibrarianID,
			Notes:         issue.Notes,
			CreatedAt:     issue.CreatedAt,
			UpdatedAt:     issue.UpdatedAt,
			ReaderName:    &issue.ReaderName,
			TicketNumber:  &issue.TicketNumber,
			BookTitle:     &issue.BookTitle,
			CopyCode:      &issue.CopyCode,
			IsOverdue:     false,
			IsActive:      true,
		}
	}

	return c.JSON(fiber.Map{"due_soon_issues": response})
}

// getMyIssues возвращает выдачи текущего читателя
func (h *Handler) getMyIssues(c *fiber.Ctx) error {
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
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get reader for issues")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	return h.getReaderIssuesInternal(c, reader.ID, true)
}

// getMyActiveIssues возвращает активные выдачи текущего читателя
func (h *Handler) getMyActiveIssues(c *fiber.Ctx) error {
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
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get reader for active issues")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	return h.getActiveIssuesByReaderInternal(c, reader.ID)
}

// getMyIssueHistory возвращает историю выдач текущего читателя
func (h *Handler) getMyIssueHistory(c *fiber.Ctx) error {
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
		log.Error().Err(err).Str("userID", userIDStr).Msg("Failed to get reader for issue history")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve reader profile.", err.Error())
	}

	return h.getReaderIssueHistoryInternal(c, reader.ID)
}

// getReaderIssues возвращает выдачи читателя по ID
func (h *Handler) getReaderIssues(c *fiber.Ctx) error {
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
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own issues.")
		}
	}

	return h.getReaderIssuesInternal(c, readerID, true)
}

// getActiveIssuesByReader возвращает активные выдачи читателя
func (h *Handler) getActiveIssuesByReader(c *fiber.Ctx) error {
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
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own issues.")
		}
	}

	return h.getActiveIssuesByReaderInternal(c, readerID)
}

// getReaderIssueHistory возвращает историю выдач читателя
func (h *Handler) getReaderIssueHistory(c *fiber.Ctx) error {
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
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own issue history.")
		}
	}

	return h.getReaderIssueHistoryInternal(c, readerID)
}

// getActiveIssuesByReaderInternal - внутренний метод для получения активных выдач читателя
func (h *Handler) getActiveIssuesByReaderInternal(c *fiber.Ctx, readerID uuid.UUID) error {
	issues, err := h.repo.GetActiveIssuesByReader(c.Context(), readerID)
	if err != nil {
		log.Error().Err(err).Str("readerID", readerID.String()).Msg("Failed to get active issues by reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve active issues.", err.Error())
	}

	response := make([]IssueResponse, len(issues))
	for i, issue := range issues {
		isOverdue := issue.OverdueDays > 0

		response[i] = IssueResponse{
			ID:            issue.ID,
			ReaderID:      issue.ReaderID,
			BookCopyID:    issue.BookCopyID,
			IssueDate:     issue.IssueDate,
			DueDate:       issue.DueDate,
			ReturnDate:    issue.ReturnDate,
			ExtendedCount: issue.ExtendedCount,
			LibrarianID:   issue.LibrarianID,
			Notes:         issue.Notes,
			CreatedAt:     issue.CreatedAt,
			UpdatedAt:     issue.UpdatedAt,
			BookTitle:     &issue.BookTitle,
			CopyCode:      &issue.CopyCode,
			OverdueDays:   &issue.OverdueDays,
			IsOverdue:     isOverdue,
			IsActive:      true,
		}
	}

	return c.JSON(fiber.Map{"active_issues": response})
}

// getReaderIssuesInternal - внутренний метод для получения выдач читателя
func (h *Handler) getReaderIssuesInternal(c *fiber.Ctx, readerID uuid.UUID, includeHistory bool) error {
	if includeHistory {
		// Получаем полную историю с пагинацией
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

		issues, err := h.repo.GetIssueHistory(c.Context(), postgres.GetIssueHistoryParams{
			ReaderID:  readerID,
			LimitVal:  limit,
			OffsetVal: offset,
		})
		if err != nil {
			log.Error().Err(err).Str("readerID", readerID.String()).Msg("Failed to get reader issue history")
			return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve issue history.", err.Error())
		}

		response := make([]IssueResponse, len(issues))
		for i, issue := range issues {
			isActive := issue.ReturnDate == nil
			var isOverdue bool
			if isActive && issue.DueDate.Before(time.Now()) {
				isOverdue = true
			}

			response[i] = IssueResponse{
				ID:            issue.ID,
				ReaderID:      issue.ReaderID,
				BookCopyID:    issue.BookCopyID,
				IssueDate:     issue.IssueDate,
				DueDate:       issue.DueDate,
				ReturnDate:    issue.ReturnDate,
				ExtendedCount: issue.ExtendedCount,
				LibrarianID:   issue.LibrarianID,
				Notes:         issue.Notes,
				CreatedAt:     issue.CreatedAt,
				UpdatedAt:     issue.UpdatedAt,
				BookTitle:     &issue.BookTitle,
				CopyCode:      &issue.CopyCode,
				IsOverdue:     isOverdue,
				IsActive:      isActive,
			}
		}

		return c.JSON(fiber.Map{"issues": response})
	} else {
		// Только активные выдачи
		return h.getActiveIssuesByReaderInternal(c, readerID)
	}
}

// getReaderIssueHistoryInternal - внутренний метод для получения истории выдач читателя
func (h *Handler) getReaderIssueHistoryInternal(c *fiber.Ctx, readerID uuid.UUID) error {
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

	issues, err := h.repo.GetIssueHistory(c.Context(), postgres.GetIssueHistoryParams{
		ReaderID:  readerID,
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		log.Error().Err(err).Str("readerID", readerID.String()).Msg("Failed to get reader issue history")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve issue history.", err.Error())
	}

	response := make([]IssueResponse, len(issues))
	for i, issue := range issues {
		isActive := issue.ReturnDate == nil
		var isOverdue bool
		if isActive && issue.DueDate.Before(time.Now()) {
			isOverdue = true
		}

		response[i] = IssueResponse{
			ID:            issue.ID,
			ReaderID:      issue.ReaderID,
			BookCopyID:    issue.BookCopyID,
			IssueDate:     issue.IssueDate,
			DueDate:       issue.DueDate,
			ReturnDate:    issue.ReturnDate,
			ExtendedCount: issue.ExtendedCount,
			LibrarianID:   issue.LibrarianID,
			Notes:         issue.Notes,
			CreatedAt:     issue.CreatedAt,
			UpdatedAt:     issue.UpdatedAt,
			BookTitle:     &issue.BookTitle,
			CopyCode:      &issue.CopyCode,
			IsOverdue:     isOverdue,
			IsActive:      isActive,
		}
	}

	return c.JSON(fiber.Map{"issue_history": response, "limit": limit, "offset": offset})
}

// getBookCopyHistory возвращает историю выдач для конкретного экземпляра книги
// func (h *Handler) getBookCopyHistory(c *fiber.Ctx) error {
// 	// Проверка роли пользователя
// 	userRole := c.Locals("userRole").(string)
// 	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
// 		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view book copy history.")
// 	}

// 	copyIDStr := c.Params("copyId")
// 	copyID, err := uuid.Parse(copyIDStr)
// 	if err != nil {
// 		return httperr.New(fiber.StatusBadRequest, "Invalid copy ID format.")
// 	}

// 	issues, err := h.repo.GetBookIssueHistory(c.Context(), copyID)
// 	if err != nil {
// 		log.Error().Err(err).Str("copyID", copyIDStr).Msg("Failed to get book copy history")
// 		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy history.", err.Error())
// 	}

// 	response := make([]IssueResponse, len(issues))
// 	for i, issue := range issues {
// 		isActive := issue.ReturnDate == nil
// 		var isOverdue bool
// 		if isActive && issue.DueDate.Before(time.Now()) {
// 			isOverdue = true
// 		}

// 		response[i] = IssueResponse{
// 			ID:            issue.ID,
// 			ReaderID:      issue.ReaderID,
// 			BookCopyID:    issue.BookCopyID,
// 			IssueDate:     issue.IssueDate,
// 			DueDate:       issue.DueDate,
// 			ReturnDate:    issue.ReturnDate,
// 			ExtendedCount: issue.ExtendedCount,
// 			LibrarianID:   issue.LibrarianID,
// 			Notes:         issue.Notes,
// 			CreatedAt:     issue.CreatedAt,
// 			UpdatedAt:     issue.UpdatedAt,
// 			ReaderName:    &issue.ReaderName,
// 			TicketNumber:  &issue.TicketNumber,
// 			IsOverdue:     isOverdue,
// 			IsActive:      isActive,
// 		}
// 	}

// 	return c.JSON(fiber.Map{"copy_history": response})
// }

// getIssueByID возвращает выдачу по ID
func (h *Handler) getIssueByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	issueID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid issue ID format.")
	}

	issue, err := h.repo.GetBookIssueByID(c.Context(), issueID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Issue not found.")
		}
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to get issue by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve issue.", err.Error())
	}

	// Проверка доступа для читателей
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != issue.ReaderID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only view your own issues.")
		}
	}

	isActive := issue.ReturnDate == nil
	var isOverdue bool
	if isActive && issue.DueDate.Before(time.Now()) {
		isOverdue = true
	}

	response := IssueResponse{
		ID:            issue.ID,
		ReaderID:      issue.ReaderID,
		BookCopyID:    issue.BookCopyID,
		IssueDate:     issue.IssueDate,
		DueDate:       issue.DueDate,
		ReturnDate:    issue.ReturnDate,
		ExtendedCount: issue.ExtendedCount,
		LibrarianID:   issue.LibrarianID,
		Notes:         issue.Notes,
		CreatedAt:     issue.CreatedAt,
		UpdatedAt:     issue.UpdatedAt,
		ReaderName:    &issue.ReaderName,
		TicketNumber:  &issue.TicketNumber,
		BookTitle:     &issue.BookTitle,
		CopyCode:      &issue.CopyCode,
		LibrarianName: issue.LibrarianName,
		IsOverdue:     isOverdue,
		IsActive:      isActive,
	}

	return c.JSON(response)
}

// createIssue создает новую выдачу книги
func (h *Handler) createIssue(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can create issues.")
	}

	var req CreateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация
	if req.DueDate.Before(time.Now()) {
		return httperr.New(fiber.StatusBadRequest, "Due date cannot be in the past.")
	}

	// Если дата выдачи не указана, используем текущую дату
	if req.IssueDate == nil {
		now := time.Now()
		req.IssueDate = &now
	}

	// Если библиотекарь не указан, берем текущего пользователя
	if req.LibrarianID == nil {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)
		req.LibrarianID = &userID
	}

	// Проверяем, что читатель существует и активен
	reader, err := h.repo.GetReaderByID(c.Context(), req.ReaderID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Reader not found.")
		}
		log.Error().Err(err).Str("readerID", req.ReaderID.String()).Msg("Failed to verify reader")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify reader.", err.Error())
	}

	if reader.IsActive != nil && !*reader.IsActive {
		return httperr.New(fiber.StatusBadRequest, "Cannot issue book to inactive reader.")
	}

	// Проверяем, что экземпляр книги существует и доступен
	bookCopy, err := h.repo.GetBookCopyByID(c.Context(), req.BookCopyID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found.")
		}
		log.Error().Err(err).Str("copyID", req.BookCopyID.String()).Msg("Failed to verify book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify book copy.", err.Error())
	}

	if bookCopy.Status.BookStatus != postgres.BookStatusAvailable {
		return httperr.New(fiber.StatusBadRequest, "Book copy is not available for issue.")
	}

	// Проверяем лимит активных выдач для читателя
	activeCount, err := h.repo.CountActiveIssuesByReader(c.Context(), req.ReaderID)
	if err != nil {
		log.Error().Err(err).Str("readerID", req.ReaderID.String()).Msg("Failed to count active issues")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify reader's active issues.", err.Error())
	}

	// Предполагаем лимит в 5 книг на читателя
	const maxActiveIssues = 5
	if activeCount >= maxActiveIssues {
		return httperr.New(fiber.StatusBadRequest, "Reader has reached the maximum limit of active book issues.")
	}

	// Создание выдачи
	issue, err := h.repo.CreateBookIssue(c.Context(), postgres.CreateBookIssueParams{
		ReaderID:    req.ReaderID,
		BookCopyID:  req.BookCopyID,
		IssueDate:   req.IssueDate,
		DueDate:     req.DueDate,
		LibrarianID: req.LibrarianID,
		Notes:       req.Notes,
	})
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid reader ID, book copy ID, or librarian ID.")
		}
		log.Error().Err(err).Msg("Failed to create book issue")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create book issue.", err.Error())
	}

	// Обновляем статус экземпляра книги на "выдан"
	_, err = h.repo.UpdateBookCopy(c.Context(), postgres.UpdateBookCopyParams{
		CopyID: req.BookCopyID,
		Status: postgres.NullBookStatus{
			BookStatus: postgres.BookStatusIssued,
			Valid:      true,
		},
	})
	if err != nil {
		log.Warn().Err(err).Str("copyID", req.BookCopyID.String()).Msg("Failed to update book copy status after issue")
	}

	response := IssueResponse{
		ID:            issue.ID,
		ReaderID:      issue.ReaderID,
		BookCopyID:    issue.BookCopyID,
		IssueDate:     issue.IssueDate,
		DueDate:       issue.DueDate,
		ReturnDate:    issue.ReturnDate,
		ExtendedCount: issue.ExtendedCount,
		LibrarianID:   issue.LibrarianID,
		Notes:         issue.Notes,
		CreatedAt:     issue.CreatedAt,
		UpdatedAt:     issue.UpdatedAt,
		IsOverdue:     false,
		IsActive:      true,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateIssue обновляет информацию о выдаче
func (h *Handler) updateIssue(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can update issues.")
	}

	idStr := c.Params("id")
	issueID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid issue ID format.")
	}

	var req UpdateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию о выдаче
	existingIssue, err := h.repo.GetBookIssueByID(c.Context(), issueID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Issue not found.")
		}
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to get issue for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve issue.", err.Error())
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateBookIssueParams{
		IssueID: issueID,
	}

	if req.DueDate != nil {
		if req.DueDate.Before(time.Now()) && existingIssue.ReturnDate == nil {
			return httperr.New(fiber.StatusBadRequest, "Due date cannot be in the past for active issues.")
		}
		updateParams.DueDate = *req.DueDate
	} else {
		updateParams.DueDate = existingIssue.DueDate
	}

	if req.ReturnDate != nil {
		updateParams.ReturnDate = req.ReturnDate
	} else {
		updateParams.ReturnDate = existingIssue.ReturnDate
	}

	if req.ExtendedCount != nil {
		updateParams.ExtendedCount = req.ExtendedCount
	} else {
		updateParams.ExtendedCount = existingIssue.ExtendedCount
	}

	if req.Notes != nil {
		updateParams.Notes = req.Notes
	} else {
		updateParams.Notes = existingIssue.Notes
	}

	// Обновление выдачи
	updatedIssue, err := h.repo.UpdateBookIssue(c.Context(), updateParams)
	if err != nil {
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to update issue")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update issue.", err.Error())
	}

	// Если устанавливается дата возврата, обновляем статус экземпляра книги
	if req.ReturnDate != nil && existingIssue.ReturnDate == nil {
		_, err = h.repo.UpdateBookCopy(c.Context(), postgres.UpdateBookCopyParams{
			CopyID: existingIssue.BookCopyID,
			Status: postgres.NullBookStatus{
				BookStatus: postgres.BookStatusAvailable,
				Valid:      true,
			},
		})
		if err != nil {
			log.Warn().Err(err).Str("copyID", existingIssue.BookCopyID.String()).Msg("Failed to update book copy status after return")
		}
	}

	isActive := updatedIssue.ReturnDate == nil
	var isOverdue bool
	if isActive && updatedIssue.DueDate.Before(time.Now()) {
		isOverdue = true
	}

	response := IssueResponse{
		ID:            updatedIssue.ID,
		ReaderID:      updatedIssue.ReaderID,
		BookCopyID:    updatedIssue.BookCopyID,
		IssueDate:     updatedIssue.IssueDate,
		DueDate:       updatedIssue.DueDate,
		ReturnDate:    updatedIssue.ReturnDate,
		ExtendedCount: updatedIssue.ExtendedCount,
		LibrarianID:   updatedIssue.LibrarianID,
		Notes:         updatedIssue.Notes,
		CreatedAt:     updatedIssue.CreatedAt,
		UpdatedAt:     updatedIssue.UpdatedAt,
		IsOverdue:     isOverdue,
		IsActive:      isActive,
	}

	return c.JSON(response)
}

// extendIssue продлевает срок выдачи
func (h *Handler) extendIssue(c *fiber.Ctx) error {
	idStr := c.Params("id")
	issueID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid issue ID format.")
	}

	var req ExtendIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация
	if req.ExtensionDays <= 0 || req.ExtensionDays > 30 {
		return httperr.New(fiber.StatusBadRequest, "Extension days must be between 1 and 30.")
	}

	// Получаем информацию о выдаче
	existingIssue, err := h.repo.GetBookIssueByID(c.Context(), issueID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Issue not found.")
		}
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to get issue for extension")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve issue.", err.Error())
	}

	// Проверка доступа для читателей
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != existingIssue.ReaderID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only extend your own issues.")
		}
	}

	// Проверяем, что выдача активна
	if existingIssue.ReturnDate != nil {
		return httperr.New(fiber.StatusBadRequest, "Cannot extend returned book issue.")
	}

	// Проверяем лимит продлений
	const maxExtensions = 3
	if existingIssue.ExtendedCount != nil && *existingIssue.ExtendedCount >= maxExtensions {
		return httperr.New(fiber.StatusBadRequest, "Maximum number of extensions reached.")
	}

	// Рассчитываем новую дату возврата
	newDueDate := existingIssue.DueDate.AddDate(0, 0, req.ExtensionDays)

	// Продлеваем выдачу
	err = h.repo.ExtendBookIssue(c.Context(), postgres.ExtendBookIssueParams{
		IssueID:    issueID,
		NewDueDate: newDueDate,
	})
	if err != nil {
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to extend issue")
		return httperr.New(fiber.StatusInternalServerError, "Failed to extend issue.", err.Error())
	}

	// Получаем обновленную информацию
	updatedIssue, err := h.repo.GetBookIssueByID(c.Context(), issueID)
	if err != nil {
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to get updated issue")
		return httperr.New(fiber.StatusInternalServerError, "Extension processed but failed to retrieve updated issue.", err.Error())
	}

	response := IssueResponse{
		ID:            updatedIssue.ID,
		ReaderID:      updatedIssue.ReaderID,
		BookCopyID:    updatedIssue.BookCopyID,
		IssueDate:     updatedIssue.IssueDate,
		DueDate:       updatedIssue.DueDate,
		ReturnDate:    updatedIssue.ReturnDate,
		ExtendedCount: updatedIssue.ExtendedCount,
		LibrarianID:   updatedIssue.LibrarianID,
		Notes:         updatedIssue.Notes,
		CreatedAt:     updatedIssue.CreatedAt,
		UpdatedAt:     updatedIssue.UpdatedAt,
		ReaderName:    &updatedIssue.ReaderName,
		TicketNumber:  &updatedIssue.TicketNumber,
		BookTitle:     &updatedIssue.BookTitle,
		CopyCode:      &updatedIssue.CopyCode,
		LibrarianName: updatedIssue.LibrarianName,
		IsOverdue:     false,
		IsActive:      true,
	}

	return c.JSON(response)
}

// returnBook обрабатывает возврат книги
func (h *Handler) returnBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	issueID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid issue ID format.")
	}

	var req ReturnBookRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем информацию о выдаче
	existingIssue, err := h.repo.GetBookIssueByID(c.Context(), issueID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Issue not found.")
		}
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to get issue for return")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve issue.", err.Error())
	}

	// Проверка доступа для читателей
	userRole := c.Locals("userRole").(string)
	if userRole == string(postgres.UserRoleReader) {
		userIDStr := c.Locals("userID").(string)
		userID, _ := uuid.Parse(userIDStr)

		reader, err := h.repo.GetReaderByUserID(c.Context(), &userID)
		if err != nil || reader.ID != existingIssue.ReaderID {
			return httperr.New(fiber.StatusForbidden, "Access denied. You can only return your own books.")
		}
	}

	// Проверяем, что книга еще не возвращена
	if existingIssue.ReturnDate != nil {
		return httperr.New(fiber.StatusBadRequest, "Book has already been returned.")
	}

	// Если дата возврата не указана, используем текущую дату
	if req.ReturnDate == nil {
		now := time.Now()
		req.ReturnDate = &now
	}

	// Валидация даты возврата
	if req.ReturnDate.Before(*existingIssue.IssueDate) {
		return httperr.New(fiber.StatusBadRequest, "Return date cannot be before issue date.")
	}

	// Обрабатываем возврат
	err = h.repo.ReturnBook(c.Context(), postgres.ReturnBookParams{
		IssueID:    issueID,
		ReturnDate: req.ReturnDate,
	})
	if err != nil {
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to return book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to process book return.", err.Error())
	}

	// Обновляем статус экземпляра книги на "доступен"
	_, err = h.repo.UpdateBookCopy(c.Context(), postgres.UpdateBookCopyParams{
		CopyID: existingIssue.BookCopyID,
		Status: postgres.NullBookStatus{
			BookStatus: postgres.BookStatusAvailable,
			Valid:      true,
		},
	})
	if err != nil {
		log.Warn().Err(err).Str("copyID", existingIssue.BookCopyID.String()).Msg("Failed to update book copy status after return")
	}

	// Получаем обновленную информацию
	updatedIssue, err := h.repo.GetBookIssueByID(c.Context(), issueID)
	if err != nil {
		log.Error().Err(err).Str("issueID", idStr).Msg("Failed to get updated issue")
		return httperr.New(fiber.StatusInternalServerError, "Return processed but failed to retrieve updated issue.", err.Error())
	}

	response := IssueResponse{
		ID:            updatedIssue.ID,
		ReaderID:      updatedIssue.ReaderID,
		BookCopyID:    updatedIssue.BookCopyID,
		IssueDate:     updatedIssue.IssueDate,
		DueDate:       updatedIssue.DueDate,
		ReturnDate:    updatedIssue.ReturnDate,
		ExtendedCount: updatedIssue.ExtendedCount,
		LibrarianID:   updatedIssue.LibrarianID,
		Notes:         updatedIssue.Notes,
		CreatedAt:     updatedIssue.CreatedAt,
		UpdatedAt:     updatedIssue.UpdatedAt,
		ReaderName:    &updatedIssue.ReaderName,
		TicketNumber:  &updatedIssue.TicketNumber,
		BookTitle:     &updatedIssue.BookTitle,
		CopyCode:      &updatedIssue.CopyCode,
		LibrarianName: updatedIssue.LibrarianName,
		IsOverdue:     false,
		IsActive:      false,
	}

	return c.JSON(response)
}

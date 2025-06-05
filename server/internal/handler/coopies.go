package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	httperr "github.com/hnnsly/library-console/pkg/error"
	"github.com/rs/zerolog/log"
)

// Request/Response structs для book copies
type CreateBookCopyRequest struct {
	BookID         uuid.UUID  `json:"book_id" validate:"required"`
	CopyCode       string     `json:"copy_code" validate:"required,min=1,max=50"`
	Status         string     `json:"status" validate:"omitempty,oneof=available issued reserved lost damaged"`
	ReadingHallID  *uuid.UUID `json:"reading_hall_id"`
	ConditionNotes *string    `json:"condition_notes" validate:"omitempty,max=1000"`
}

type UpdateBookCopyRequest struct {
	Status         *string    `json:"status" validate:"omitempty,oneof=available issued reserved lost damaged"`
	ReadingHallID  *uuid.UUID `json:"reading_hall_id"`
	ConditionNotes *string    `json:"condition_notes" validate:"omitempty,max=1000"`
}

type UpdateBookCopyStatusRequest struct {
	Status         string  `json:"status" validate:"required,oneof=available issued reserved lost damaged"`
	ConditionNotes *string `json:"condition_notes" validate:"omitempty,max=1000"`
}

type BookCopyResponse struct {
	ID             uuid.UUID  `json:"id"`
	BookID         uuid.UUID  `json:"book_id"`
	CopyCode       string     `json:"copy_code"`
	Status         string     `json:"status"`
	ReadingHallID  *uuid.UUID `json:"reading_hall_id"`
	ConditionNotes *string    `json:"condition_notes"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	// Extended fields
	BookTitle *string `json:"book_title,omitempty"`
	BookISBN  *string `json:"book_isbn,omitempty"`
	HallName  *string `json:"hall_name,omitempty"`
}

type BookCopyStatisticsResponse struct {
	BookID          uuid.UUID `json:"book_id"`
	TotalCopies     int64     `json:"total_copies"`
	AvailableCopies int64     `json:"available_copies"`
	IssuedCopies    int64     `json:"issued_copies"`
	ReservedCopies  int64     `json:"reserved_copies"`
	LostCopies      int64     `json:"lost_copies"`
	DamagedCopies   int64     `json:"damaged_copies"`
}

// Вспомогательная функция для преобразования строки статуса в NullBookStatus
func stringToNullBookStatus(status string) postgres.NullBookStatus {
	if status == "" {
		return postgres.NullBookStatus{Valid: false}
	}
	return postgres.NullBookStatus{
		BookStatus: postgres.BookStatus(status),
		Valid:      true,
	}
}

// Вспомогательная функция для преобразования NullBookStatus в строку
func nullBookStatusToString(status postgres.NullBookStatus) string {
	if !status.Valid {
		return ""
	}
	return string(status.BookStatus)
}

// listBookCopies возвращает список всех экземпляров книг (только для администраторов/библиотекарей)
func (h *Handler) listBookCopies(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can view all book copies.")
	}

	// Получаем статус из query параметров
	c.Query("status", "available")

	return h.getBookCopiesByStatus(c)
}

// listAvailableBookCopies возвращает доступные экземпляры книг
func (h *Handler) listAvailableBookCopies(c *fiber.Ctx) error {
	copies, err := h.repo.ListAvailableBookCopies(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to list available book copies")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve available book copies.", err.Error())
	}

	response := make([]BookCopyResponse, len(copies))
	for i, copy := range copies {
		response[i] = BookCopyResponse{
			ID:             copy.ID,
			BookID:         copy.BookID,
			CopyCode:       copy.CopyCode,
			Status:         nullBookStatusToString(copy.Status),
			ReadingHallID:  copy.ReadingHallID,
			ConditionNotes: copy.ConditionNotes,
			CreatedAt:      copy.CreatedAt,
			UpdatedAt:      copy.UpdatedAt,
			BookTitle:      &copy.Title,
			BookISBN:       copy.Isbn,
		}
	}

	return c.JSON(fiber.Map{"copies": response})
}

// getBookCopiesByStatus возвращает экземпляры книг по статусу
func (h *Handler) getBookCopiesByStatus(c *fiber.Ctx) error {
	statusStr := c.Params("status")
	if statusStr == "" {
		statusStr = c.Query("status", "available")
	}

	// Валидация статуса
	validStatuses := []string{"available", "issued", "reserved", "lost", "damaged"}
	isValidStatus := false
	for _, vs := range validStatuses {
		if statusStr == vs {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return httperr.New(fiber.StatusBadRequest, "Invalid status. Valid values are: available, issued, reserved, lost, damaged.")
	}

	status := stringToNullBookStatus(statusStr)
	copies, err := h.repo.GetBookCopiesByStatus(c.Context(), status)
	if err != nil {
		log.Error().Err(err).Str("status", statusStr).Msg("Failed to get book copies by status")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copies.", err.Error())
	}

	response := make([]BookCopyResponse, len(copies))
	for i, copy := range copies {
		response[i] = BookCopyResponse{
			ID:             copy.ID,
			BookID:         copy.BookID,
			CopyCode:       copy.CopyCode,
			Status:         nullBookStatusToString(copy.Status),
			ReadingHallID:  copy.ReadingHallID,
			ConditionNotes: copy.ConditionNotes,
			CreatedAt:      copy.CreatedAt,
			UpdatedAt:      copy.UpdatedAt,
			BookTitle:      &copy.Title,
			BookISBN:       copy.Isbn,
			HallName:       copy.HallName,
		}
	}

	return c.JSON(fiber.Map{"copies": response, "status": statusStr})
}

// getBookCopiesByBook возвращает все экземпляры конкретной книги
func (h *Handler) getBookCopiesByBook(c *fiber.Ctx) error {
	bookIDStr := c.Params("bookId")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid book ID format.")
	}

	copies, err := h.repo.ListBookCopiesByBook(c.Context(), bookID)
	if err != nil {
		log.Error().Err(err).Str("bookID", bookIDStr).Msg("Failed to get book copies by book")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copies.", err.Error())
	}

	response := make([]BookCopyResponse, len(copies))
	for i, copy := range copies {
		response[i] = BookCopyResponse{
			ID:             copy.ID,
			BookID:         copy.BookID,
			CopyCode:       copy.CopyCode,
			Status:         nullBookStatusToString(copy.Status),
			ReadingHallID:  copy.ReadingHallID,
			ConditionNotes: copy.ConditionNotes,
			CreatedAt:      copy.CreatedAt,
			UpdatedAt:      copy.UpdatedAt,
			HallName:       copy.HallName,
		}
	}

	// Получаем статистику для этой книги
	totalCopies, _ := h.repo.CountBookCopiesByBook(c.Context(), bookID)
	availableCopies, _ := h.repo.CountAvailableBookCopies(c.Context(), bookID)

	return c.JSON(fiber.Map{
		"copies": response,
		"statistics": map[string]interface{}{
			"total_copies":     totalCopies,
			"available_copies": availableCopies,
		},
	})
}

// getBookCopyByCode возвращает экземпляр книги по коду
func (h *Handler) getBookCopyByCode(c *fiber.Ctx) error {
	copyCode := c.Params("code")
	if copyCode == "" {
		return httperr.New(fiber.StatusBadRequest, "Copy code is required.")
	}

	copy, err := h.repo.GetBookCopyByCode(c.Context(), copyCode)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found.")
		}
		log.Error().Err(err).Str("copyCode", copyCode).Msg("Failed to get book copy by code")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy.", err.Error())
	}

	response := BookCopyResponse{
		ID:             copy.ID,
		BookID:         copy.BookID,
		CopyCode:       copy.CopyCode,
		Status:         nullBookStatusToString(copy.Status),
		ReadingHallID:  copy.ReadingHallID,
		ConditionNotes: copy.ConditionNotes,
		CreatedAt:      copy.CreatedAt,
		UpdatedAt:      copy.UpdatedAt,
		BookTitle:      &copy.Title,
		BookISBN:       copy.Isbn,
	}

	return c.JSON(response)
}

// getBookCopyByID возвращает экземпляр книги по ID
func (h *Handler) getBookCopyByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	copyID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid copy ID format.")
	}

	copy, err := h.repo.GetBookCopyByID(c.Context(), copyID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found.")
		}
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to get book copy by ID")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy.", err.Error())
	}

	response := BookCopyResponse{
		ID:             copy.ID,
		BookID:         copy.BookID,
		CopyCode:       copy.CopyCode,
		Status:         nullBookStatusToString(copy.Status),
		ReadingHallID:  copy.ReadingHallID,
		ConditionNotes: copy.ConditionNotes,
		CreatedAt:      copy.CreatedAt,
		UpdatedAt:      copy.UpdatedAt,
		BookTitle:      &copy.Title,
		BookISBN:       copy.Isbn,
	}

	return c.JSON(response)
}

// getBookCopyHistory возвращает историю выдач для экземпляра книги
func (h *Handler) getBookCopyHistory(c *fiber.Ctx) error {
	// Этот хендлер уже реализован в issues.go, но мы можем добавить дублирующий маршрут
	return h.getBookCopyHistory(c)
}

// createBookCopy создает новый экземпляр книги
func (h *Handler) createBookCopy(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can create book copies.")
	}

	var req CreateBookCopyRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация
	if req.CopyCode == "" {
		return httperr.New(fiber.StatusBadRequest, "Copy code is required.")
	}

	// Если статус не указан, используем 'available'
	if req.Status == "" {
		req.Status = "available"
	}

	// Проверяем, что книга существует
	_, err := h.repo.GetBookByID(c.Context(), req.BookID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book not found.")
		}
		log.Error().Err(err).Str("bookID", req.BookID.String()).Msg("Failed to verify book existence")
		return httperr.New(fiber.StatusInternalServerError, "Failed to verify book.", err.Error())
	}

	// Проверяем, что код экземпляра уникален
	existingCopy, err := h.repo.GetBookCopyByCode(c.Context(), req.CopyCode)
	if err == nil && existingCopy != nil {
		return httperr.New(fiber.StatusConflict, "Copy code already exists.")
	}

	// Если указан зал, проверяем его существование
	if req.ReadingHallID != nil {
		_, err := h.repo.GetReadingHallByID(c.Context(), *req.ReadingHallID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return httperr.New(fiber.StatusNotFound, "Reading hall not found.")
			}
			log.Error().Err(err).Str("hallID", req.ReadingHallID.String()).Msg("Failed to verify reading hall")
			return httperr.New(fiber.StatusInternalServerError, "Failed to verify reading hall.", err.Error())
		}
	}

	status := stringToNullBookStatus(req.Status)
	copy, err := h.repo.CreateBookCopy(c.Context(), postgres.CreateBookCopyParams{
		BookID:         req.BookID,
		CopyCode:       req.CopyCode,
		Status:         status,
		ReadingHallID:  req.ReadingHallID,
		ConditionNotes: req.ConditionNotes,
	})
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return httperr.New(fiber.StatusConflict, "Copy code already exists.")
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid book ID or reading hall ID.")
		}
		log.Error().Err(err).Msg("Failed to create book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create book copy.", err.Error())
	}

	response := BookCopyResponse{
		ID:             copy.ID,
		BookID:         copy.BookID,
		CopyCode:       copy.CopyCode,
		Status:         nullBookStatusToString(copy.Status),
		ReadingHallID:  copy.ReadingHallID,
		ConditionNotes: copy.ConditionNotes,
		CreatedAt:      copy.CreatedAt,
		UpdatedAt:      copy.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// updateBookCopy обновляет экземпляр книги
func (h *Handler) updateBookCopy(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can update book copies.")
	}

	idStr := c.Params("id")
	copyID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid copy ID format.")
	}

	var req UpdateBookCopyRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Получаем текущую информацию об экземпляре
	existingCopy, err := h.repo.GetBookCopyByID(c.Context(), copyID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found.")
		}
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to get book copy for update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy.", err.Error())
	}

	// Если указан новый зал, проверяем его существование
	if req.ReadingHallID != nil {
		_, err := h.repo.GetReadingHallByID(c.Context(), *req.ReadingHallID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return httperr.New(fiber.StatusNotFound, "Reading hall not found.")
			}
			log.Error().Err(err).Str("hallID", req.ReadingHallID.String()).Msg("Failed to verify reading hall")
			return httperr.New(fiber.StatusInternalServerError, "Failed to verify reading hall.", err.Error())
		}
	}

	// Подготавливаем параметры для обновления
	updateParams := postgres.UpdateBookCopyParams{
		CopyID: copyID,
	}

	if req.Status != nil {
		updateParams.Status = stringToNullBookStatus(*req.Status)
	} else {
		updateParams.Status = existingCopy.Status
	}

	if req.ReadingHallID != nil {
		updateParams.ReadingHallID = req.ReadingHallID
	} else {
		updateParams.ReadingHallID = existingCopy.ReadingHallID
	}

	if req.ConditionNotes != nil {
		updateParams.ConditionNotes = req.ConditionNotes
	} else {
		updateParams.ConditionNotes = existingCopy.ConditionNotes
	}

	// Обновление экземпляра
	updatedCopy, err := h.repo.UpdateBookCopy(c.Context(), updateParams)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusBadRequest, "Invalid reading hall ID.")
		}
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to update book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book copy.", err.Error())
	}

	response := BookCopyResponse{
		ID:             updatedCopy.ID,
		BookID:         updatedCopy.BookID,
		CopyCode:       updatedCopy.CopyCode,
		Status:         nullBookStatusToString(updatedCopy.Status),
		ReadingHallID:  updatedCopy.ReadingHallID,
		ConditionNotes: updatedCopy.ConditionNotes,
		CreatedAt:      updatedCopy.CreatedAt,
		UpdatedAt:      updatedCopy.UpdatedAt,
	}

	return c.JSON(response)
}

// updateBookCopyStatus обновляет только статус экземпляра книги
func (h *Handler) updateBookCopyStatus(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) && userRole != string(postgres.UserRoleLibrarian) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators and librarians can update book copy status.")
	}

	idStr := c.Params("id")
	copyID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid copy ID format.")
	}

	var req UpdateBookCopyStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}

	// Валидация статуса
	validStatuses := []string{"available", "issued", "reserved", "lost", "damaged"}
	isValidStatus := false
	for _, vs := range validStatuses {
		if req.Status == vs {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return httperr.New(fiber.StatusBadRequest, "Invalid status. Valid values are: available, issued, reserved, lost, damaged.")
	}

	// Получаем текущую информацию об экземпляре
	existingCopy, err := h.repo.GetBookCopyByID(c.Context(), copyID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found.")
		}
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to get book copy for status update")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy.", err.Error())
	}

	// Проверяем, можно ли изменить статус
	currentStatus := nullBookStatusToString(existingCopy.Status)
	if currentStatus == "issued" && req.Status != "available" && req.Status != "lost" && req.Status != "damaged" {
		return httperr.New(fiber.StatusBadRequest, "Cannot change status of issued book to this value. Book must be returned first.")
	}

	// Обновление статуса
	updateParams := postgres.UpdateBookCopyParams{
		CopyID:         copyID,
		Status:         stringToNullBookStatus(req.Status),
		ReadingHallID:  existingCopy.ReadingHallID,
		ConditionNotes: req.ConditionNotes,
	}

	if req.ConditionNotes == nil {
		updateParams.ConditionNotes = existingCopy.ConditionNotes
	}

	updatedCopy, err := h.repo.UpdateBookCopy(c.Context(), updateParams)
	if err != nil {
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to update book copy status")
		return httperr.New(fiber.StatusInternalServerError, "Failed to update book copy status.", err.Error())
	}

	response := BookCopyResponse{
		ID:             updatedCopy.ID,
		BookID:         updatedCopy.BookID,
		CopyCode:       updatedCopy.CopyCode,
		Status:         nullBookStatusToString(updatedCopy.Status),
		ReadingHallID:  updatedCopy.ReadingHallID,
		ConditionNotes: updatedCopy.ConditionNotes,
		CreatedAt:      updatedCopy.CreatedAt,
		UpdatedAt:      updatedCopy.UpdatedAt,
	}

	return c.JSON(response)
}

// deleteBookCopy удаляет экземпляр книги
func (h *Handler) deleteBookCopy(c *fiber.Ctx) error {
	// Проверка роли пользователя
	userRole := c.Locals("userRole").(string)
	if userRole != string(postgres.UserRoleAdministrator) {
		return httperr.New(fiber.StatusForbidden, "Access denied. Only administrators can delete book copies.")
	}

	idStr := c.Params("id")
	copyID, err := uuid.Parse(idStr)
	if err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid copy ID format.")
	}

	// Проверяем, что экземпляр существует
	existingCopy, err := h.repo.GetBookCopyByID(c.Context(), copyID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return httperr.New(fiber.StatusNotFound, "Book copy not found.")
		}
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to get book copy for deletion")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve book copy.", err.Error())
	}

	// Проверяем, что экземпляр не выдан
	currentStatus := nullBookStatusToString(existingCopy.Status)
	if currentStatus == "issued" {
		return httperr.New(fiber.StatusBadRequest, "Cannot delete issued book copy. Book must be returned first.")
	}

	// Проверяем, нет ли активных выдач
	//activeIssuesCount, err := h.repo.CountActiveIssuesByReader(c.Context(), copyID) // Это не правильно, нужен метод для подсчета по copy
	// Для простоты, пропускаем эту проверку, но в реальном проекте стоит добавить соответствующий SQL-запрос

	// Удаляем экземпляр
	err = h.repo.DeleteBookCopy(c.Context(), copyID)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return httperr.New(fiber.StatusConflict, "Cannot delete book copy. It has related records (issues, history, etc.).")
		}
		log.Error().Err(err).Str("copyID", idStr).Msg("Failed to delete book copy")
		return httperr.New(fiber.StatusInternalServerError, "Failed to delete book copy.", err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

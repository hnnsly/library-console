package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

// createDailyStatistics создает или обновляет ежедневную статистику
func (h *Handler) createDailyStatistics(c *fiber.Ctx) error {
	err := h.repo.CreateDailyStatistics(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create daily statistics")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create daily statistics.")
	}

	return c.JSON(fiber.Map{"message": "Daily statistics successfully created/updated"})
}

// getLoanStatusStatistics получает статистику по статусам выдач за период
func (h *Handler) getLoanStatusStatistics(c *fiber.Ctx) error {
	daysBackStr := c.Query("days_back", "30") // default 30 days
	daysBack, err := strconv.Atoi(daysBackStr)
	if err != nil || daysBack <= 0 {
		return httperr.New(fiber.StatusBadRequest, "Invalid days_back parameter.")
	}

	// TODO: Validate days_back is reasonable (max 365 days)

	if daysBack > 365 {
		return httperr.New(fiber.StatusBadRequest, "days_back cannot exceed 365 days.")
	}

	fromDate := time.Now().AddDate(0, 0, -daysBack)
	statistics, err := h.repo.GetLoanStatusStatistics(c.Context(), fromDate)
	if err != nil {
		log.Error().Err(err).Int("daysBack", daysBack).Msg("Failed to get loan status statistics")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve loan status statistics.")
	}

	if statistics == nil {
		statistics = []*postgres.GetLoanStatusStatisticsRow{}
	}

	return c.JSON(statistics)
}

// getMonthlyReport получает месячный отчет за последние 12 месяцев
func (h *Handler) getMonthlyReport(c *fiber.Ctx) error {
	report, err := h.repo.GetMonthlyReport(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get monthly report")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve monthly report.")
	}

	if report == nil {
		report = []*postgres.GetMonthlyReportRow{}
	}

	return c.JSON(report)
}

// getYearlyReportByCategory получает годовой отчет по категориям
func (h *Handler) getYearlyReportByCategory(c *fiber.Ctx) error {
	report, err := h.repo.GetYearlyReportByCategory(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get yearly report by category")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve yearly report by category.")
	}

	if report == nil {
		report = []*postgres.GetYearlyReportByCategoryRow{}
	}

	return c.JSON(report)
}

// getInventoryReport получает отчет по инвентаризации
func (h *Handler) getInventoryReport(c *fiber.Ctx) error {
	report, err := h.repo.GetInventoryReport(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get inventory report")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve inventory report.")
	}

	if report == nil {
		report = []*postgres.GetInventoryReportRow{}
	}

	return c.JSON(report)
}

// getLibraryOverview получает общий обзор библиотеки (комбинированная статистика)
func (h *Handler) getLibraryOverview(c *fiber.Ctx) error {
	// Получаем различные статистики для создания общего обзора

	// Статистика по выдачам за последние 30 дней
	fromDate := time.Now().AddDate(0, 0, -30)
	loanStats, err := h.repo.GetLoanStatusStatistics(c.Context(), fromDate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get loan statistics for overview")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve library overview.")
	}

	// Общее количество читателей
	totalReaders, err := h.repo.GetReadersCount(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get readers count for overview")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve library overview.")
	}

	// Просроченные книги
	overdueBooks, err := h.repo.GetOverdueBooks(c.Context(), 100)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get overdue books for overview")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve library overview.")
	}

	// Книги к возврату сегодня
	booksDueToday, err := h.repo.GetBooksDueToday(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get books due today for overview")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve library overview.")
	}

	// Неоплаченные штрафы
	unpaidFines, err := h.repo.GetUnpaidFines(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get unpaid fines for overview")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve library overview.")
	}

	// Статистика по залам
	hallStats, err := h.repo.GetHallStatistics(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get hall statistics for overview")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve library overview.")
	}

	overview := map[string]interface{}{
		"total_readers":       totalReaders,
		"loan_statistics":     loanStats,
		"overdue_books_count": len(overdueBooks),
		"overdue_books":       overdueBooks,
		"books_due_today":     booksDueToday,
		"unpaid_fines_count":  len(unpaidFines),
		"unpaid_fines":        unpaidFines,
		"hall_statistics":     hallStats,
		"generated_at":        time.Now(),
	}

	return c.JSON(overview)
}

// getDashboardStats получает основные метрики для дашборда
func (h *Handler) getDashboardStats(c *fiber.Ctx) error {
	// Получаем ключевые метрики для дашборда

	// Общее количество читателей
	totalReaders, err := h.repo.GetReadersCount(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get readers count for dashboard")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve dashboard statistics.")
	}

	// Количество активных выдач
	fromDate := time.Now().AddDate(0, 0, -1) // за последний день
	loanStats, err := h.repo.GetLoanStatusStatistics(c.Context(), fromDate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get loan statistics for dashboard")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve dashboard statistics.")
	}

	// Просроченные книги
	overdueBooks, err := h.repo.GetOverdueBooks(c.Context(), 1000)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get overdue books for dashboard")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve dashboard statistics.")
	}

	// Книги к возврату сегодня
	booksDueToday, err := h.repo.GetBooksDueToday(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get books due today for dashboard")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve dashboard statistics.")
	}

	// Должники
	debtors, err := h.repo.GetDebtorReaders(c.Context())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get debtors for dashboard")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve dashboard statistics.")
	}

	// Активные читатели за последние 30 дней
	activeReaders, err := h.repo.GetActiveReaders(c.Context(), 50)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get active readers for dashboard")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve dashboard statistics.")
	}

	dashboardStats := map[string]interface{}{
		"total_readers":         totalReaders,
		"active_readers_count":  len(activeReaders),
		"overdue_books_count":   len(overdueBooks),
		"books_due_today_count": len(booksDueToday),
		"debtors_count":         len(debtors),
		"loan_statistics":       loanStats,
		"updated_at":            time.Now(),
	}

	return c.JSON(dashboardStats)
}

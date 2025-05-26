package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/config"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/middleware"
	"github.com/hnnsly/library-console/internal/repository"
	"github.com/rs/zerolog/log"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
	maxLimit      = 100
)

// Handler holds the dependencies for the library API.
type Handler struct {
	repo *repository.LibraryRepository
	cfg  config.LibraryServiceConfig
}

// NewHandler creates a new library API handler.
func NewHandler(repo *repository.LibraryRepository, cfg config.LibraryServiceConfig) *Handler {
	return &Handler{
		repo: repo,
		cfg:  cfg,
	}
}

// Router sets up the routes for the library API.
func (h *Handler) Router() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: httperr.GlobalErrorHandler,
	})

	api := app.Group("/api")

	// Books routes
	bookRoutes := api.Group("/books")
	bookRoutes.Post("/", middleware.RequireAuthHeader(), h.createBook)
	bookRoutes.Get("/", h.listBooks)
	bookRoutes.Get("/search", h.searchBooks)
	bookRoutes.Get("/available", h.getAvailableBooks)
	bookRoutes.Get("/popular", h.getPopularBooks)
	bookRoutes.Get("/top-rated", h.getTopRatedBooks)
	bookRoutes.Get("/single-copy", h.getBooksWithSingleCopy)
	bookRoutes.Get("/:id", h.getBookByID)
	bookRoutes.Put("/:id/copies", middleware.RequireAuthHeader(), h.updateBookCopies)
	bookRoutes.Delete("/:id/writeoff", middleware.RequireAuthHeader(), h.writeOffBook)
	bookRoutes.Get("/:id/queue", h.getBookQueue)
	bookRoutes.Get("/:id/loans", h.getActiveLoansByBook)

	// Readers routes
	readerRoutes := api.Group("/readers")
	readerRoutes.Post("/", middleware.RequireAuthHeader(), h.createReader)
	readerRoutes.Get("/", h.listReaders)
	readerRoutes.Get("/search", h.searchReaders)
	readerRoutes.Get("/active", h.getActiveReaders)
	readerRoutes.Get("/debtors", h.getDebtorReaders)
	readerRoutes.Get("/:id", h.getReaderByID)
	readerRoutes.Put("/:id", middleware.RequireAuthHeader(), h.updateReader)
	readerRoutes.Put("/:id/status", middleware.RequireAuthHeader(), h.updateReaderStatus)
	readerRoutes.Get("/:id/loans", h.getReaderCurrentLoans)
	readerRoutes.Get("/:id/history", h.getReaderLoanHistory)
	readerRoutes.Get("/:id/statistics", h.getReaderStatistics)
	readerRoutes.Get("/:id/fines", h.getReaderFines)
	readerRoutes.Get("/:id/reservations", h.getReaderReservations)
	readerRoutes.Get("/:id/favorites", h.getReaderFavoriteCategories)

	// Loans routes
	loanRoutes := api.Group("/loans")
	loanRoutes.Post("/", middleware.RequireAuthHeader(), h.createLoan)
	loanRoutes.Get("/overdue", h.getOverdueBooks)
	loanRoutes.Get("/due-today", h.getBooksDueToday)
	loanRoutes.Get("/:id", h.getLoanByID)
	loanRoutes.Put("/:id/return", middleware.RequireAuthHeader(), h.returnBook)
	loanRoutes.Put("/:id/renew", middleware.RequireAuthHeader(), h.renewLoan)
	loanRoutes.Put("/:id/lost", middleware.RequireAuthHeader(), h.markLoanAsLost)
	loanRoutes.Get("/:id/renewals", h.getRenewalsForLoan)
	loanRoutes.Post("/:id/check-eligibility", h.checkLoanEligibility)

	// Fines routes
	fineRoutes := api.Group("/fines")
	fineRoutes.Post("/", middleware.RequireAuthHeader(), h.createFine)
	fineRoutes.Get("/unpaid", h.getUnpaidFines)
	fineRoutes.Put("/:id/pay", middleware.RequireAuthHeader(), h.payFine)
	fineRoutes.Put("/:id/waive", middleware.RequireAuthHeader(), h.waiveFine)

	// Reservations routes
	reservationRoutes := api.Group("/reservations")
	reservationRoutes.Post("/", middleware.RequireAuthHeader(), h.createReservation)
	reservationRoutes.Get("/expired", h.getExpiredReservations)
	reservationRoutes.Put("/:id/fulfill", middleware.RequireAuthHeader(), h.fulfillReservation)
	reservationRoutes.Put("/:id/cancel", middleware.RequireAuthHeader(), h.cancelReservation)

	// Halls routes
	hallRoutes := api.Group("/halls")
	hallRoutes.Get("/", h.getAllHalls)
	hallRoutes.Get("/:id", h.getHallByID)
	hallRoutes.Put("/:id/occupancy", middleware.RequireAuthHeader(), h.updateHallOccupancy)
	hallRoutes.Get("/:id/statistics", h.getHallStatistics)
	hallRoutes.Get("/:id/books/by-author", h.getBooksByAuthorInHall)

	// Categories routes
	categoryRoutes := api.Group("/categories")
	categoryRoutes.Post("/", middleware.RequireAuthHeader(), h.createCategory)
	categoryRoutes.Get("/", h.getAllCategories)
	categoryRoutes.Get("/statistics", h.getCategoryStatistics)

	// Librarians routes
	librarianRoutes := api.Group("/librarians")
	librarianRoutes.Post("/", middleware.RequireAuthHeader(), h.createLibrarian)
	librarianRoutes.Get("/", h.getAllLibrarians)
	librarianRoutes.Get("/:id", h.getLibrarianByID)
	librarianRoutes.Put("/:id", middleware.RequireAuthHeader(), h.updateLibrarian)
	librarianRoutes.Put("/:id/deactivate", middleware.RequireAuthHeader(), h.deactivateLibrarian)

	// Search routes
	searchRoutes := api.Group("/search")
	searchRoutes.Get("/global", h.globalSearch)
	searchRoutes.Get("/advanced", h.advancedBookSearch)

	// Statistics routes
	statsRoutes := api.Group("/statistics")
	statsRoutes.Get("/loans", h.getLoanStatusStatistics)
	statsRoutes.Get("/monthly", h.getMonthlyReport)
	statsRoutes.Get("/yearly", h.getYearlyReportByCategory)
	statsRoutes.Get("/inventory", h.getInventoryReport)
	statsRoutes.Post("/daily", middleware.RequireAuthHeader(), h.createDailyStatistics)

	// Operations routes
	operationRoutes := api.Group("/operations")
	operationRoutes.Get("/logs", h.getOperationLogs)
	operationRoutes.Get("/recent", h.getRecentOperations)

	// Renewals routes
	renewalRoutes := api.Group("/renewals")
	renewalRoutes.Get("/by-date", h.getRenewalsByDate)
	renewalRoutes.Get("/most-renewed", h.getMostRenewedBooks)

	app.Get("/health", h.healthCheck)

	return app
}

func (h *Handler) healthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"status": "healthy"})
}

// --- Helper Functions ---
func parsePagination(c *fiber.Ctx) (limit, offset int32) {
	limitStr := c.Query("limit", strconv.Itoa(defaultLimit))
	offsetStr := c.Query("offset", strconv.Itoa(defaultOffset))

	limitInt, err := strconv.Atoi(limitStr)
	if err != nil || limitInt <= 0 {
		limitInt = defaultLimit
	}
	if limitInt > maxLimit {
		limitInt = maxLimit
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil || offsetInt < 0 {
		offsetInt = defaultOffset
	}

	return int32(limitInt), int32(offsetInt)
}

func parseID(c *fiber.Ctx, paramName string) (int64, error) {
	idStr := c.Params(paramName)
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, httperr.New(fiber.StatusBadRequest, fmt.Sprintf("Invalid %s format. Must be an integer.", paramName))
	}
	return id, nil
}

func parseOptionalID(c *fiber.Ctx, queryName string) int32 {
	idStr := c.Query(queryName, "0")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0
	}
	return int32(id)
}

func parseYear(c *fiber.Ctx, queryName string) int32 {
	yearStr := c.Query(queryName, "0")
	year, err := strconv.ParseInt(yearStr, 10, 32)
	if err != nil {
		return 0
	}
	return int32(year)
}

func validateBodyAndParse(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		log.Warn().Err(err).Msg("Failed to parse request body")
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.", err.Error())
	}
	return nil
}

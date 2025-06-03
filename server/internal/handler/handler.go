package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hnnsly/library-console/internal/auth"
	"github.com/hnnsly/library-console/internal/config"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/logger"
	"github.com/hnnsly/library-console/internal/middleware"
	"github.com/hnnsly/library-console/internal/repository"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/hnnsly/library-console/internal/repository/redis"
	"github.com/rs/zerolog/log"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
	maxLimit      = 100
)

// Handler holds the dependencies for the library API.
type Handler struct {
	repo           *repository.LibraryRepository
	cfg            config.LibraryServiceConfig
	authHandler    *AuthHandler
	sessionManager *auth.SessionManager
}

// NewHandler creates a new library API handler.
func NewHandler(repo *repository.LibraryRepository, cfg config.LibraryServiceConfig, pgQueries *postgres.Queries, redisClient *redis.Redis) *Handler {
	sessionManager := auth.NewSessionManager(redisClient)
	authHandler := NewAuthHandler(pgQueries, sessionManager, 12*time.Hour)

	return &Handler{
		repo:           repo,
		cfg:            cfg,
		authHandler:    authHandler,
		sessionManager: sessionManager,
	}
}

// Router sets up the routes for the library API.
func (h *Handler) Router() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: httperr.GlobalErrorHandler,
	})

	app.Use(logger.RequestLogger())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://lab.somerka.ru", // Убедитесь что это правильный домен
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,HEAD,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token,Cookie",
		AllowCredentials: true,         // ✅ Правильно
		ExposeHeaders:    "Set-Cookie", // Добавить для cookies
		MaxAge:           86400,
	}))

	// Health check route (no auth required)
	app.Get("/health", h.healthCheck)

	// Auth routes (public)
	h.setupAuthRoutes(app)

	// API routes (protected)
	api := app.Group("/api")

	// Apply authentication middleware to all API routes
	api.Use(middleware.AuthRequired(h.sessionManager, 12*time.Hour))

	// Setup all protected routes
	h.setupBooksRoutes(api)
	h.setupReadersRoutes(api)
	h.setupLoansRoutes(api)
	h.setupFinesRoutes(api)
	h.setupReservationsRoutes(api)
	h.setupHallsRoutes(api)
	h.setupCategoriesRoutes(api)
	h.setupLibrariansRoutes(api)
	h.setupUsersRoutes(api)
	h.setupSearchRoutes(api)
	h.setupStatisticsRoutes(api)
	h.setupOperationsRoutes(api)
	h.setupRenewalsRoutes(api)

	return app
}

// setupAuthRoutes настраивает роуты авторизации
func (h *Handler) setupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	// Публичные роуты
	auth.Post("/login", h.authHandler.login)

	// Защищенные роуты авторизации
	authProtected := auth.Use(middleware.AuthRequired(h.sessionManager, 12*time.Hour))
	authProtected.Post("/logout", h.authHandler.logout)
	authProtected.Get("/me", h.authHandler.getCurrentUser)
	authProtected.Post("/change-password", h.authHandler.changePassword)
}

// setupBooksRoutes настраивает роуты для работы с книгами
func (h *Handler) setupBooksRoutes(api fiber.Router) {
	books := api.Group("/books")

	// Роуты доступные всем авторизованным пользователям
	books.Get("/", h.getAvailableBooks)
	books.Get("/search", h.searchBooks)
	books.Get("/available", h.getAvailableBooks)
	books.Get("/popular", h.getPopularBooks)
	books.Get("/top-rated", h.getTopRatedBooks)
	books.Get("/single-copy", h.getBooksWithSingleCopy)
	books.Get("/:id", h.getBookByID)
	books.Get("/code/:code", h.getBookByCode)
	books.Get("/:id/queue", h.getBookQueue)
	books.Get("/:id/loans", h.getActiveLoansByBook)

	// Роуты для библиотекарей и администраторов
	librarianOnly := books.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/", h.createBook)
	librarianOnly.Put("/:id/availability", h.updateBookAvailability)
	librarianOnly.Put("/:id/copies", h.updateBookCopies)
	librarianOnly.Delete("/:id/writeoff", h.writeOffBook)
	librarianOnly.Post("/advanced-search", h.advancedSearchBooks)
}

// setupReadersRoutes настраивает роуты для работы с читателями
func (h *Handler) setupReadersRoutes(api fiber.Router) {
	readers := api.Group("/readers")

	// Роуты доступные всем авторизованным пользователям
	readers.Get("/", h.getAllReaders)
	readers.Get("/active", h.getActiveReaders)
	readers.Get("/debtors", h.getDebtorReaders)
	readers.Get("/count", h.getReadersCount)
	readers.Get("/:id", h.getReaderByID)
	readers.Get("/ticket/:ticket", h.getReaderByTicket)
	readers.Get("/:id/loans", h.getReaderCurrentLoans)
	readers.Get("/:id/statistics", h.getReaderStatistics)
	readers.Get("/:id/fines", h.getReaderFines)
	readers.Get("/:id/reservations", h.getReaderReservations)
	readers.Get("/:id/favorites", h.getReaderFavoriteCategories)
	readers.Post("/search", h.searchReadersByName)
	readers.Post("/:id/history", h.getReaderLoanHistory)

	// Роуты для библиотекарей и администраторов
	librarianOnly := readers.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/", h.createReader)
	librarianOnly.Put("/:id", h.updateReader)
	librarianOnly.Put("/:id/status", h.updateReaderStatus)
	librarianOnly.Put("/:id/debt", h.updateReaderDebt)
}

// setupLoansRoutes настраивает роуты для работы с выдачами
func (h *Handler) setupLoansRoutes(api fiber.Router) {
	loans := api.Group("/loans")

	// Роуты доступные всем авторизованным пользователям
	loans.Get("/overdue", h.getOverdueBooks)
	loans.Get("/due-today", h.getBooksDueToday)
	loans.Get("/:id", h.getLoanByID)
	loans.Get("/:id/renewals", h.getRenewalsForLoan)
	loans.Get("/check-eligibility", h.checkLoanEligibility)

	// Роуты для библиотекарей и администраторов
	librarianOnly := loans.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/", h.createLoan)
	librarianOnly.Put("/:id/return", h.returnBook)
	librarianOnly.Put("/:id/renew", h.renewLoan)
	librarianOnly.Put("/:id/lost", h.markLoanAsLost)
	librarianOnly.Post("/renewals", h.createRenewal)
}

// setupFinesRoutes настраивает роуты для работы со штрафами
func (h *Handler) setupFinesRoutes(api fiber.Router) {
	fines := api.Group("/fines")

	// Роуты доступные всем авторизованным пользователям
	fines.Get("/unpaid", h.getUnpaidFines)

	// Роуты для библиотекарей и администраторов
	librarianOnly := fines.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/", h.createFine)
	librarianOnly.Post("/calculate", h.calculateOverdueFine)
	librarianOnly.Put("/:id/pay", h.payFine)
	librarianOnly.Put("/:id/waive", h.waiveFine)
}

// setupReservationsRoutes настраивает роуты для работы с бронированиями
func (h *Handler) setupReservationsRoutes(api fiber.Router) {
	reservations := api.Group("/reservations")

	// Роуты доступные всем авторизованным пользователям
	reservations.Get("/expired", h.getExpiredReservations)

	// Роуты для библиотекарей и администраторов
	librarianOnly := reservations.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/", h.createReservation)
	librarianOnly.Put("/:id/fulfill", h.fulfillReservation)
	librarianOnly.Put("/:id/cancel", h.cancelReservation)
}

// setupHallsRoutes настраивает роуты для работы с залами
func (h *Handler) setupHallsRoutes(api fiber.Router) {
	halls := api.Group("/halls")

	// Роуты доступные всем авторизованным пользователям
	halls.Get("/", h.getAllHalls)
	halls.Get("/:id", h.getHallByID)
	halls.Get("/statistics", h.getHallStatistics)
	halls.Get("/books/by-author", h.getBooksByAuthorInHall)

	// Роуты для библиотекарей и администраторов
	librarianOnly := halls.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Put("/:id/occupancy", h.updateHallOccupancy)
}

// setupCategoriesRoutes настраивает роуты для работы с категориями
func (h *Handler) setupCategoriesRoutes(api fiber.Router) {
	categories := api.Group("/categories")

	// Роуты доступные всем авторизованным пользователям
	categories.Get("/", h.getAllCategories)
	categories.Get("/statistics", h.getCategoryStatistics)

	// Роуты для администраторов
	adminOnly := categories.Use(middleware.RequireRole("admin", "super_admin"))
	adminOnly.Post("/", h.createCategory)
}

// setupLibrariansRoutes настраивает роуты для работы с библиотекарями
func (h *Handler) setupLibrariansRoutes(api fiber.Router) {
	librarians := api.Group("/librarians")

	// Роуты доступные всем авторизованным пользователям
	librarians.Get("/", h.getAllLibrarians)
	librarians.Get("/:id", h.getLibrarianByID)
	librarians.Get("/employee/:employee_id", h.getLibrarianByEmployeeID)

	// Роуты для администраторов
	adminOnly := librarians.Use(middleware.RequireRole("admin", "super_admin"))
	adminOnly.Post("/", h.createLibrarian)
	adminOnly.Put("/:id", h.updateLibrarian)
	adminOnly.Put("/:id/deactivate", h.deactivateLibrarian)
}

// setupUsersRoutes настраивает роуты для управления пользователями
func (h *Handler) setupUsersRoutes(api fiber.Router) {
	// users := api.Group("/users")

	// // Роуты для администраторов
	// adminOnly := users.Use(middleware.RequireRole("admin", "super_admin"))
	// adminOnly.Post("/", h.createUser)
	// adminOnly.Post("/list", h.getAllUsers)
	// adminOnly.Get("/:id", h.getUserByID)
	// adminOnly.Put("/:id", h.updateUser)
	// adminOnly.Put("/:id/role", h.updateUserRole)
	// adminOnly.Put("/:id/deactivate", h.deactivateUser)
	// adminOnly.Put("/:id/activate", h.activateUser)
	// adminOnly.Get("/role/:role", h.getUsersByRole)

	// // Роуты только для super_admin
	// superAdminOnly := users.Use(middleware.RequireRole("super_admin"))
	// superAdminOnly.Delete("/:id", h.deleteUser)
}

// setupSearchRoutes настраивает роуты для поиска
func (h *Handler) setupSearchRoutes(api fiber.Router) {
	search := api.Group("/search")

	// Роуты доступные всем авторизованным пользователям
	search.Get("/global", h.globalSearch)
	search.Post("/books/advanced", h.advancedBookSearch)
}

// setupStatisticsRoutes настраивает роуты для статистики
func (h *Handler) setupStatisticsRoutes(api fiber.Router) {
	stats := api.Group("/statistics")

	// Роуты доступные всем авторизованным пользователям
	stats.Get("/loans", h.getLoanStatusStatistics)
	stats.Get("/monthly", h.getMonthlyReport)
	stats.Get("/yearly", h.getYearlyReportByCategory)
	stats.Get("/inventory", h.getInventoryReport)
	stats.Get("/overview", h.getLibraryOverview)
	stats.Get("/dashboard", h.getDashboardStats)

	// Роуты для библиотекарей и администраторов
	librarianOnly := stats.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/daily", h.createDailyStatistics)
}

// setupOperationsRoutes настраивает роуты для журнала операций
func (h *Handler) setupOperationsRoutes(api fiber.Router) {
	operations := api.Group("/operations")

	// Роуты доступные всем авторизованным пользователям
	operations.Post("/logs", h.getOperationLogs)
	operations.Get("/recent", h.getRecentOperations)

	// Роуты для библиотекарей и администраторов
	librarianOnly := operations.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/", h.createOperationLog)
}

// setupRenewalsRoutes настраивает роуты для продлений
func (h *Handler) setupRenewalsRoutes(api fiber.Router) {
	renewals := api.Group("/renewals")

	// Роуты доступные всем авторизованным пользователям
	renewals.Post("/by-date", h.getRenewalsByDate)
	renewals.Get("/most-renewed", h.getMostRenewedBooks)

	// Роуты для библиотекарей и администраторов
	librarianOnly := renewals.Use(middleware.RequireRole("librarian", "admin", "super_admin"))
	librarianOnly.Post("/", h.createRenewalRecord)
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
	id, err := strconv.ParseInt(idStr, 10, 64) // Изменено на 64 бита
	if err != nil {
		return 0, httperr.New(fiber.StatusBadRequest, fmt.Sprintf("Invalid %s format. Must be an integer.", paramName))
	}
	return id, nil
}

func parseOptionalID(c *fiber.Ctx, queryName string) int64 {
	idStr := c.Query(queryName, "0")
	id, err := strconv.ParseInt(idStr, 10, 64) // Изменено на 64 бита
	if err != nil {
		return 0
	}
	return id
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

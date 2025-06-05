package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hnnsly/library-console/internal/config"
	"github.com/hnnsly/library-console/internal/logger"
	"github.com/hnnsly/library-console/internal/middleware"
	"github.com/hnnsly/library-console/internal/repository"
	httperr "github.com/hnnsly/library-console/pkg/error"

	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/hnnsly/library-console/internal/repository/redis"
)

type Handler struct {
	repo *repository.LibraryRepository
	cfg  *config.LibraryServiceConfig
}

func NewHandler(pg *postgres.Queries, rdb *redis.Redis, cfg *config.LibraryServiceConfig) *Handler {
	return &Handler{repo: repository.New(pg, rdb), cfg: cfg}
}

func (h *Handler) Router() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: httperr.GlobalErrorHandler,
	})

	app.Use(logger.RequestLogger())
	authMW := middleware.NewAuthMiddleware(h.repo, 24*time.Hour)

	// Health check
	app.Get("/health", h.healthCheck)

	// API routes
	api := app.Group("/api/library")

	// Auth routes (no auth required)
	authRoutes := api.Group("/auth")
	authRoutes.Post("/login", h.login)
	authRoutes.Post("/register", h.register)
	authRoutes.Post("/logout", authMW, h.logout)

	// User management routes (auth required)
	userRoutes := api.Group("/users", authMW)
	userRoutes.Get("/", h.listUsers)
	userRoutes.Get("/me", h.getCurrentUser)
	userRoutes.Get("/:id", h.getUserByID)
	userRoutes.Put("/:id", h.updateUser)
	userRoutes.Delete("/:id", h.deactivateUser)
	userRoutes.Post("/", h.createUser)
	userRoutes.Put("/:id/password", h.updateUserPassword)
	userRoutes.Get("/role/:role", h.getUsersByRole)

	// Reader management routes (auth required)
	readerRoutes := api.Group("/readers", authMW)
	readerRoutes.Get("/", h.listReaders)
	readerRoutes.Get("/me", h.getCurrentReader)
	readerRoutes.Get("/search", h.searchReaders)
	readerRoutes.Get("/ticket/:ticket", h.getReaderByTicket)
	readerRoutes.Get("/hall/:hallId", h.getReadersByHall)
	readerRoutes.Get("/:id", h.getReaderByID)
	readerRoutes.Post("/", h.createReader)
	readerRoutes.Put("/:id", h.updateReader)
	readerRoutes.Delete("/:id", h.deactivateReader)

	// Reading halls management routes (auth required)
	hallRoutes := api.Group("/halls", authMW)
	hallRoutes.Get("/", h.listReadingHalls)
	hallRoutes.Get("/statistics", h.getHallStatistics)
	hallRoutes.Get("/:id", h.getReadingHallByID)
	hallRoutes.Post("/", h.createReadingHall)
	hallRoutes.Put("/:id", h.updateReadingHall)
	hallRoutes.Delete("/:id", h.deleteReadingHall)
	hallRoutes.Put("/:id/occupancy", h.updateHallOccupancy)

	// Authors management routes (auth required)
	authorRoutes := api.Group("/authors", authMW)
	authorRoutes.Get("/", h.listAuthors)
	authorRoutes.Get("/search", h.searchAuthors)
	authorRoutes.Get("/:id", h.getAuthorByID)
	authorRoutes.Get("/:id/books", h.getAuthorBooks)
	authorRoutes.Post("/", h.createAuthor)
	authorRoutes.Put("/:id", h.updateAuthor)
	authorRoutes.Delete("/:id", h.deleteAuthor)

	// Fines management routes (auth required)
	fineRoutes := api.Group("/fines", authMW)
	fineRoutes.Get("/", h.listFines)
	fineRoutes.Get("/unpaid", h.getUnpaidFines)
	fineRoutes.Get("/statistics", h.getFineStatistics)
	fineRoutes.Get("/my", h.getMyFines)
	fineRoutes.Get("/reader/:readerId", h.getFinesByReader)
	fineRoutes.Get("/reader/:readerId/unpaid", h.getUnpaidFinesByReader)
	fineRoutes.Get("/reader/:readerId/debt", h.getReaderDebt)
	fineRoutes.Get("/:id", h.getFineByID)
	fineRoutes.Post("/", h.createFine)
	fineRoutes.Put("/:id", h.updateFine)
	fineRoutes.Post("/:id/pay", h.payFine)
	fineRoutes.Delete("/:id", h.deleteFine)

	// Books management routes (auth required)
	bookRoutes := api.Group("/books", authMW)
	bookRoutes.Get("/", h.listBooks)
	bookRoutes.Get("/search", h.searchBooks)
	bookRoutes.Get("/top-rated", h.getTopRatedBooks)
	bookRoutes.Get("/author/:authorId", h.getBooksByAuthor)
	bookRoutes.Get("/isbn/:isbn", h.getBookByISBN)
	bookRoutes.Get("/:id", h.getBookByID)
	bookRoutes.Get("/:id/details", h.getBookWithDetails)
	bookRoutes.Get("/:id/authors", h.getBookAuthors)
	bookRoutes.Post("/", h.createBook)
	bookRoutes.Put("/:id", h.updateBook)
	bookRoutes.Put("/:id/availability", h.updateBookAvailability)
	bookRoutes.Post("/:id/authors", h.addBookAuthor)
	bookRoutes.Delete("/:id/authors/:authorId", h.removeBookAuthor)
	bookRoutes.Delete("/:id/authors", h.removeAllBookAuthors)
	bookRoutes.Delete("/:id", h.deleteBook)

	// Book copies management routes (auth required)
	copyRoutes := api.Group("/copies", authMW)
	copyRoutes.Get("/", h.listBookCopies)
	copyRoutes.Get("/available", h.listAvailableBookCopies)
	copyRoutes.Get("/status/:status", h.getBookCopiesByStatus)
	copyRoutes.Get("/book/:bookId", h.getBookCopiesByBook)
	copyRoutes.Get("/code/:code", h.getBookCopyByCode)
	copyRoutes.Get("/:id", h.getBookCopyByID)
	copyRoutes.Get("/:id/history", h.getBookCopyHistory)
	copyRoutes.Post("/", h.createBookCopy)
	copyRoutes.Put("/:id", h.updateBookCopy)
	copyRoutes.Put("/:id/status", h.updateBookCopyStatus)
	copyRoutes.Delete("/:id", h.deleteBookCopy)

	// Book ratings management routes (auth required)
	ratingRoutes := api.Group("/ratings", authMW)
	ratingRoutes.Get("/my", h.getMyRatings)
	ratingRoutes.Get("/reader/:readerId", h.getReaderRatings)
	ratingRoutes.Get("/book/:bookId", h.getBookRatings)
	ratingRoutes.Get("/book/:bookId/average", h.getBookAverageRating)
	ratingRoutes.Get("/book/:bookId/my", h.getMyBookRating)
	ratingRoutes.Get("/:id", h.getRatingByID)
	ratingRoutes.Post("/", h.createRating)
	ratingRoutes.Put("/:id", h.updateRating)
	ratingRoutes.Delete("/:id", h.deleteRating)

	// Book issues management routes (auth required)
	issueRoutes := api.Group("/issues", authMW)
	issueRoutes.Get("/", h.listIssues)
	issueRoutes.Get("/active", h.getActiveIssues)
	issueRoutes.Get("/overdue", h.getOverdueIssues)
	issueRoutes.Get("/due-soon", h.getIssuesDueSoon)
	issueRoutes.Get("/my", h.getMyIssues)
	issueRoutes.Get("/my/active", h.getMyActiveIssues)
	issueRoutes.Get("/my/history", h.getMyIssueHistory)
	issueRoutes.Get("/reader/:readerId", h.getReaderIssues)
	issueRoutes.Get("/reader/:readerId/active", h.getActiveIssuesByReader)
	issueRoutes.Get("/reader/:readerId/history", h.getReaderIssueHistory)
	issueRoutes.Get("/copy/:copyId/history", h.getBookCopyHistory)
	issueRoutes.Get("/:id", h.getIssueByID)
	issueRoutes.Post("/", h.createIssue)
	issueRoutes.Put("/:id", h.updateIssue)
	issueRoutes.Post("/:id/extend", h.extendIssue)
	issueRoutes.Post("/:id/return", h.returnBook)

	return app
}

func (h *Handler) healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "healthy"})
}

package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hnnsly/library-console/internal/config"
	"github.com/hnnsly/library-console/internal/middleware"
	"github.com/hnnsly/library-console/internal/repository"
	httperr "github.com/hnnsly/library-console/pkg/error"
)

type Handler struct {
	repo *repository.LibraryRepository
	cfg  *config.LibraryServiceConfig
}

func NewHandler(repo *repository.LibraryRepository, cfg *config.LibraryServiceConfig) *Handler {
	return &Handler{
		repo: repo,
		cfg:  cfg,
	}
}

func (h *Handler) Router() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: httperr.GlobalErrorHandler,
	})

	// Middleware
	app.Use(recover.New())

	// Health check
	app.Get("/health", h.healthCheck)

	// API routes
	api := app.Group("/api/library")

	// Auth middleware for protected routes
	authMiddleware := middleware.NewAuthMiddleware(h.repo, 24*time.Hour)

	authGroup := api.Group("/auth")
	authGroup.Post("/login", h.login)
	authGroup.Post("/logout", h.logout)
	authGroup.Get("/me", authMiddleware, h.me)

	// Books
	booksGroup := api.Group("/books")
	booksGroup.Get("/", h.getAllBooks)
	booksGroup.Get("/search", h.searchBooks)
	booksGroup.Get("/:id", h.getBookById)
	booksGroup.Post("/", authMiddleware, h.createBook)
	booksGroup.Put("/:id", authMiddleware, h.updateBook)
	booksGroup.Get("/:id/authors", h.getBookAuthors)
	booksGroup.Post("/:id/authors", authMiddleware, h.addBookAuthor)
	booksGroup.Delete("/:id/authors/:authorId", authMiddleware, h.removeBookAuthor)

	// Book copies
	copiesGroup := api.Group("/copies")
	copiesGroup.Get("/book/:bookId", h.getBookCopiesByBookId)
	copiesGroup.Get("/hall/:hallId", h.getBookCopiesByHall)
	copiesGroup.Get("/code/:copyCode", h.getBookCopyByCode)
	copiesGroup.Get("/:id", h.getBookCopyById)
	copiesGroup.Post("/", authMiddleware, h.createBookCopy)
	copiesGroup.Put("/:id/status", authMiddleware, h.updateBookCopyStatus)

	// Authors
	authorsGroup := api.Group("/authors")
	authorsGroup.Get("/", h.getAllAuthors)
	authorsGroup.Get("/search", h.searchAuthors)
	authorsGroup.Get("/:id", h.getAuthorById)
	authorsGroup.Post("/", authMiddleware, h.createAuthor)
	authorsGroup.Get("/:id/books", h.getAuthorBooks)

	// Readers
	readersGroup := api.Group("/readers")
	readersGroup.Get("/", authMiddleware, h.getActiveReaders)
	readersGroup.Get("/search", authMiddleware, h.searchReaders)
	readersGroup.Get("/:id", authMiddleware, h.getReaderById)
	readersGroup.Get("/ticket/:ticketNumber", authMiddleware, h.getReaderByTicketNumber)
	readersGroup.Post("/", authMiddleware, h.createReader)
	readersGroup.Put("/:id", authMiddleware, h.updateReader)
	readersGroup.Delete("/:id", authMiddleware, h.deactivateReader)
	readersGroup.Get("/:id/books", authMiddleware, h.getReaderActiveBooks)
	readersGroup.Get("/:id/fines", authMiddleware, h.getReaderFines)
	readersGroup.Get("/:id/visits", authMiddleware, h.getReaderVisitHistory)

	// Book issues
	issuesGroup := api.Group("/issues")
	issuesGroup.Get("/", authMiddleware, h.getBooksToReturn)
	issuesGroup.Get("/overdue", authMiddleware, h.getOverdueBooks)
	issuesGroup.Get("/recent", authMiddleware, h.getRecentBookOperations)
	issuesGroup.Post("/", authMiddleware, h.issueBook)
	issuesGroup.Post("/return", authMiddleware, h.returnBook)

	// Reading halls
	hallsGroup := api.Group("/halls")
	hallsGroup.Get("/", h.getAllReadingHalls)
	hallsGroup.Get("/dashboard", authMiddleware, h.getHallsDashboard)
	hallsGroup.Get("/:id", h.getReadingHallById)
	hallsGroup.Post("/", authMiddleware, h.createReadingHall)
	hallsGroup.Put("/:id", authMiddleware, h.updateReadingHall)
	hallsGroup.Get("/:id/visits/stats/daily", authMiddleware, h.getDailyVisitStats)
	hallsGroup.Get("/:id/visits/stats/hourly", authMiddleware, h.getHourlyVisitStats)

	// Hall visits
	visitsGroup := api.Group("/visits")
	visitsGroup.Get("/recent", authMiddleware, h.getRecentHallVisits)
	visitsGroup.Post("/entry", authMiddleware, h.registerHallEntry)
	visitsGroup.Post("/exit", authMiddleware, h.registerHallExit)

	// Fines
	finesGroup := api.Group("/fines")
	finesGroup.Get("/unpaid", authMiddleware, h.getUnpaidFines)
	finesGroup.Post("/", authMiddleware, h.createFine)
	finesGroup.Post("/:id/pay", authMiddleware, h.payFine)

	// Users
	usersGroup := api.Group("/users")
	usersGroup.Get("/", authMiddleware, h.getAllUsers)
	usersGroup.Get("/:id", authMiddleware, h.getUserById)
	usersGroup.Post("/", authMiddleware, h.createUser)
	usersGroup.Put("/:id", authMiddleware, h.updateUser)
	usersGroup.Delete("/:id", authMiddleware, h.deactivateUser)

	return app
}

func (h *Handler) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "library",
	})
}

package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type GlobalSearchRequest struct {
	SearchTerm string `json:"search_term"`
}

type AdvancedBookSearchRequest struct {
	TitleFilter    string `json:"title_filter"`
	AuthorFilter   string `json:"author_filter"`
	YearFilter     int    `json:"year_filter"`
	CategoryFilter int    `json:"category_filter"`
	HallFilter     int    `json:"hall_filter"`
	AvailableOnly  bool   `json:"available_only"`
	SortBy         string `json:"sort_by"`
	PageOffset     int32  `json:"page_offset"`
	PageLimit      int32  `json:"page_limit"`
}

// globalSearch выполняет глобальный поиск по книгам и читателям
func (h *Handler) globalSearch(c *fiber.Ctx) error {
	searchTerm := c.Query("q")
	if searchTerm == "" {
		return httperr.New(fiber.StatusBadRequest, "Search term is required.")
	}

	// TODO: Validate search_term min length 2, max length 100
	// TODO: Sanitize search term to prevent injection

	if len(strings.TrimSpace(searchTerm)) < 2 {
		return httperr.New(fiber.StatusBadRequest, "Search term must be at least 2 characters long.")
	}

	results, err := h.repo.GlobalSearch(c.Context(), searchTerm)
	if err != nil {
		log.Error().Err(err).Str("searchTerm", searchTerm).Msg("Failed to perform global search")
		return httperr.New(fiber.StatusInternalServerError, "Failed to perform search.")
	}

	if results == nil {
		results = []*postgres.GlobalSearchRow{}
	}

	return c.JSON(results)
}

// advancedBookSearch выполняет расширенный поиск книг (альтернативная реализация)
func (h *Handler) advancedBookSearch(c *fiber.Ctx) error {
	req := new(AdvancedBookSearchRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate at least one search parameter is provided
	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0
	// TODO: Validate year_filter if provided
	// TODO: Validate sort_by is one of: title, author, year, popularity
	// TODO: Validate hall_filter and category_filter if provided

	if req.PageLimit == 0 {
		req.PageLimit = 20 // default limit
	}

	if req.SortBy == "" {
		req.SortBy = "title" // default sort
	}

	// Используем существующий метод AdvancedSearchBooks из репозитория
	params := postgres.AdvancedSearchBooksParams{
		TitleFilter:    req.TitleFilter,
		AuthorFilter:   req.AuthorFilter,
		YearFilter:     req.YearFilter,
		CategoryFilter: req.CategoryFilter,
		HallFilter:     req.HallFilter,
		AvailableOnly:  req.AvailableOnly,
		SortBy:         req.SortBy,
		PageOffset:     req.PageOffset,
		PageLimit:      req.PageLimit,
	}

	books, err := h.repo.AdvancedSearchBooks(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform advanced book search")
		return httperr.New(fiber.StatusInternalServerError, "Failed to search books.")
	}

	if books == nil {
		books = []*postgres.AdvancedSearchBooksRow{}
	}

	return c.JSON(books)
}

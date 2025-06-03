package handler

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	httperr "github.com/hnnsly/library-console/internal/error"
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type CreateOperationLogRequest struct {
	OperationType string                 `json:"operation_type"`
	EntityType    string                 `json:"entity_type"`
	EntityID      int                    `json:"entity_id"`
	LibrarianID   *int                   `json:"librarian_id"`
	Details       map[string]interface{} `json:"details"`
	Description   *string                `json:"description"`
}

type GetOperationLogsRequest struct {
	EntityType string `json:"entity_type"`
	EntityID   int    `json:"entity_id"`
	PageOffset int32  `json:"page_offset"`
	PageLimit  int32  `json:"page_limit"`
}

// createOperationLog создает запись в журнале операций
func (h *Handler) createOperationLog(c *fiber.Ctx) error {
	req := new(CreateOperationLogRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate required fields: operation_type, entity_type, entity_id
	// TODO: Validate operation_type is one of: create, update, delete, loan, return, renew
	// TODO: Validate entity_type is one of: book, reader, loan, fine, reservation
	// TODO: Validate entity_id > 0
	// TODO: Validate librarian_id if provided

	// Преобразуем details в JSON bytes
	var detailsBytes []byte
	if req.Details != nil {
		var err error
		detailsBytes, err = json.Marshal(req.Details)
		if err != nil {
			return httperr.New(fiber.StatusBadRequest, "Invalid details format.")
		}
	}

	params := postgres.CreateOperationLogParams{
		OperationType: req.OperationType,
		EntityType:    req.EntityType,
		EntityID:      req.EntityID,
		LibrarianID:   req.LibrarianID,
		Details:       detailsBytes,
		Description:   req.Description,
	}

	err := h.repo.CreateOperationLog(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create operation log")
		return httperr.New(fiber.StatusInternalServerError, "Failed to create operation log.")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Operation log successfully created"})
}

// getOperationLogs получает журнал операций с фильтрацией и пагинацией
func (h *Handler) getOperationLogs(c *fiber.Ctx) error {
	req := new(GetOperationLogsRequest)
	if err := c.BodyParser(req); err != nil {
		return httperr.New(fiber.StatusBadRequest, "Invalid request body.")
	}

	// TODO: Validate page_limit > 0 and <= 100
	// TODO: Validate page_offset >= 0
	// TODO: Validate entity_type if provided
	// TODO: Validate entity_id if provided

	if req.PageLimit == 0 {
		req.PageLimit = 50 // default limit
	}

	params := postgres.GetOperationLogsParams{
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		PageOffset: req.PageOffset,
		PageLimit:  req.PageLimit,
	}

	logs, err := h.repo.GetOperationLogs(c.Context(), params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get operation logs")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve operation logs.")
	}

	if logs == nil {
		logs = []*postgres.GetOperationLogsRow{}
	}

	return c.JSON(logs)
}

// getRecentOperations получает последние операции
func (h *Handler) getRecentOperations(c *fiber.Ctx) error {
	limit := int32(20) // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	// TODO: Validate limit > 0 and <= 100

	operations, err := h.repo.GetRecentOperations(c.Context(), limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent operations")
		return httperr.New(fiber.StatusInternalServerError, "Failed to retrieve recent operations.")
	}

	if operations == nil {
		operations = []*postgres.GetRecentOperationsRow{}
	}

	return c.JSON(operations)
}

package error

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ErrorResponse is the standard JSON error response structure.
type ErrorResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"` // Optional additional details
}

// Error makes ErrorResponse satisfy the error interface.
func (e *ErrorResponse) Error() string {
	return e.Message
}

// New creates a new ErrorResponse instance.
// The first detail in details, if provided, will be used as e.Details.
func New(statusCode int, message string, details ...interface{}) *ErrorResponse {
	er := &ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
	}
	if len(details) > 0 {
		er.Details = details[0]
	}
	return er
}

// GlobalErrorHandler is a centralized error handler for the Fiber application.
// It ensures that errors are returned in a consistent JSON format.
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	// Default error code and message
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	var details interface{}

	var e *ErrorResponse
	if errors.As(err, &e) {
		// If the error is already an ErrorResponse, use its properties
		code = e.StatusCode
		message = e.Message
		details = e.Details
	} else {
		// For other types of errors, log them (if not a simple fiber.ErrNotFound or similar)
		// and return a generic 500 error.
		// You might want to check for specific common errors like fiber.ErrNotFound
		// and customize their response.
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			message = fiberErr.Message
		} else {
			// Log unexpected errors
			log.Error().Err(err).Str("path", c.Path()).Msg("Unhandled error occurred")
		}
	}

	// Ensure status code is set on the response context
	c.Status(code)
	return c.JSON(ErrorResponse{
		StatusCode: code,
		Message:    message,
		Details:    details,
	})
}

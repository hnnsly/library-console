package httperr

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// HttpError представляет собой пользовательскую структуру ошибки для HTTP-ответов.
type HttpError struct {
	Code    int    `json:"-"` // HTTP status code
	Message string `json:"message"`
	Details string `json:"details,omitempty"` // Дополнительные детали для отладки или информации
}

func (e *HttpError) Error() string {
	return e.Message
}

// New создает новую HttpError.
func New(code int, message string, details ...string) *HttpError {
	he := &HttpError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		he.Details = details[0]
	}
	return he
}

// GlobalErrorHandler является глобальным обработчиком ошибок для Fiber.
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	var httpErr *HttpError
	if errors.As(err, &httpErr) {
		return c.Status(httpErr.Code).JSON(httpErr)
	}

	// Обработка ошибок валидации Fiber (если используется c.BodyParser со структурой с тегами validate)
	// Fiber может возвращать fiber.Error, который содержит информацию о валидации.
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		return c.Status(fiberError.Code).JSON(HttpError{
			Message: fiberError.Message,
		})
	}

	// Логирование непредвиденной ошибки
	log.Error().Err(err).Str("path", c.Path()).Msg("Unhandled error")

	// Общий ответ для непредвиденных ошибок
	return c.Status(fiber.StatusInternalServerError).JSON(HttpError{
		Message: "Internal Server Error",
		Details: "An unexpected error occurred.",
	})
}

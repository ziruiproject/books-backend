package http

import ivalidator "starter/internal/core/validator"

type WebResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func SuccessResponse[T any](data T, message string) WebResponse[T] {
	if message == "" {
		message = "Success"
	}
	return WebResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) WebResponse[any] {
	if message == "" {
		message = "An error occurred"
	}
	return WebResponse[any]{
		Success: false,
		Message: message,
		Data:    nil,
	}
}

func ValidationResponse(data ivalidator.ValidationErrors) WebResponse[any] {
	return WebResponse[any]{
		Success: false,
		Message: "Validation Error",
		Data:    data.Errors,
	}
}

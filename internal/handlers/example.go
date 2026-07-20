package handlers

import (
	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/services"
	"vessl.dev/vessl/internal/utils"
)

type ExampleHandler struct {
	service *services.ExampleService
}

func NewExampleHandler(s *services.ExampleService) *ExampleHandler {
	return &ExampleHandler{service: s}
}

func (h *ExampleHandler) List(c echo.Context) error {
	examples, err := h.service.ListExamples()
	if err != nil {
		return utils.Error(c, 500, "failed to list examples")
	}
	return utils.Success(c, "Available examples", examples)
}

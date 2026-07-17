package utils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	Data          any     `json:"data,omitempty"`
	Path          string  `json:"path,omitempty"`
	ExecutionTime float64 `json:"executionTime,omitempty"`
}

type PaginatedData struct {
	Records    any `json:"records"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
}

func Success(c echo.Context, message string, data interface{}) error {
	start := time.Now()
	if v := c.Get("startTime"); v != nil {
		if s, ok := v.(time.Time); ok {
			start = s
		}
	}
	execTime := time.Since(start).Seconds()

	return c.JSON(http.StatusOK, APIResponse{
		Status:        "success",
		Message:       message,
		Data:          data,
		Path:          c.Request().URL.Path,
		ExecutionTime: execTime,
	})
}

func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Path:    c.Request().URL.Path,
	})
}

func Accepted(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusAccepted, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Path:    c.Request().URL.Path,
	})
}

func Error(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: message,
		Path:    c.Request().URL.Path,
	})
}

func Paginated(c echo.Context, message string, records interface{}, total, page, limit int) error {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return Success(c, message, PaginatedData{
		Records:    records,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	})
}

package handlers

import (
	"github.com/labstack/echo/v4"

	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type StorageHandler struct {
	storageService *services.StorageService
}

func NewStorageHandler(s *services.StorageService) *StorageHandler {
	return &StorageHandler{storageService: s}
}

func (h *StorageHandler) ListStorage(c echo.Context) error {
	storages, err := h.storageService.ListStorage(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if storages == nil {
		storages = []*models.Storage{}
	}
	return c.JSON(http.StatusOK, storages)
}

func (h *StorageHandler) CreateStorage(c echo.Context) error {
	var st models.Storage
	if err := c.Bind(&st); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	created, err := h.storageService.CreateStorageWithDefaults(c.Request().Context(), &st)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, created)
}

func (h *StorageHandler) GetStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing storage id parameter"})
	}
	st, err := h.storageService.GetStorage(c.Request().Context(), id)
	if err != nil || st == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "storage record not found"})
	}
	return c.JSON(http.StatusOK, st)
}

func (h *StorageHandler) DeleteStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing storage id parameter"})
	}
	if err := h.storageService.DeleteStorage(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *StorageHandler) StartStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing storage id parameter"})
	}
	st, err := h.storageService.StartStorage(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, st)
}

func (h *StorageHandler) StopStorage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing storage id parameter"})
	}
	if err := h.storageService.StopStorage(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "stopped"})
}

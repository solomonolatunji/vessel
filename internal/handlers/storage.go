package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessl.dev/vessl/internal/models"
	"vessl.dev/vessl/internal/services"
)

type StorageHandler struct {
	storageService *services.StorageService
}

func NewStorageHandler(s *services.StorageService) *StorageHandler {
	return &StorageHandler{storageService: s}
}

// @Summary ListStorage endpoint
// @Description ListStorage endpoint
// @Tags Storage
// @Accept json
// @Produce json
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

// @Summary CreateStorage endpoint
// @Description CreateStorage endpoint
// @Tags Storage
// @Accept json
// @Produce json
// @Param request body models.Storage true "Payload"
// @Router /storage [post]
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

// @Summary GetStorage endpoint
// @Description GetStorage endpoint
// @Tags Storage
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /storage/{id} [get]
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

// @Summary DeleteStorage endpoint
// @Description DeleteStorage endpoint
// @Tags Storage
// @Accept json
// @Produce json
// @Param id path string true "id"
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

// @Summary StartStorage endpoint
// @Description StartStorage endpoint
// @Tags Storage
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /storage/{id}/start [post]
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

// @Summary StopStorage endpoint
// @Description StopStorage endpoint
// @Tags Storage
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Router /storage/{id}/stop [post]
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

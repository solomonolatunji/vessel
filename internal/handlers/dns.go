package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/models"
	"codedock.dev/codedock/internal/services"
	"codedock.dev/codedock/internal/utils"
)

type DNSHandler struct {
	dnsService *services.DNSService
}

func NewDNSHandler(dnsService *services.DNSService) *DNSHandler {
	return &DNSHandler{dnsService: dnsService}
}

func (h *DNSHandler) Create(c echo.Context) error {
	var req models.CreateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.CreateRecord(c.Request().Context(), &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Created(c, "DNS record created successfully", record)
}

func (h *DNSHandler) List(c echo.Context) error {
	domain := c.QueryParam("domain")

	records, err := h.dnsService.ListByDomain(c.Request().Context(), domain)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS records fetched successfully", records)
}

func (h *DNSHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	var req models.UpdateDNSRecordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	record, err := h.dnsService.UpdateRecord(c.Request().Context(), id, &req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "DNS record updated successfully", record)
}

func (h *DNSHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing record id")
	}

	if err := h.dnsService.DeleteRecord(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

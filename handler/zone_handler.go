package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/maftunhasansiam/SpotSync/dto"
	"github.com/maftunhasansiam/SpotSync/service"
)

type ZoneHandler struct {
	zoneService service.ZoneService
	validate    *validator.Validate
}

func NewZoneHandler(zoneService service.ZoneService, validate *validator.Validate) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
		validate:    validate,
	}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", formatValidationErrors(err)))
	}

	zone, err := h.zoneService.CreateZone(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create parking zone", err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse("Parking zone created successfully", zone))
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to retrieve parking zones", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse("Parking zones retrieved successfully", zones))
}

func (h *ZoneHandler) GetZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid zone ID", nil))
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if err.Error() == "parking zone not found" {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to retrieve parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse("Parking zone retrieved successfully", zone))
}

func (h *ZoneHandler) UpdateZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid zone ID", nil))
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", formatValidationErrors(err)))
	}

	zone, err := h.zoneService.UpdateZone(uint(id), &req)
	if err != nil {
		if err.Error() == "parking zone not found" {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to update parking zone", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse("Parking zone updated successfully", zone))
}

func (h *ZoneHandler) DeleteZone(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid zone ID", nil))
	}

	if err := h.zoneService.DeleteZone(uint(id)); err != nil {
		if err.Error() == "parking zone not found" {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to delete parking zone", err.Error()))
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse("Parking zone deleted successfully", nil))
}



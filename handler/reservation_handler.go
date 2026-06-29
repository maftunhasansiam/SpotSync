package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/maftunhasansiam/SpotSync/dto"
	"github.com/maftunhasansiam/SpotSync/repository"
	"github.com/maftunhasansiam/SpotSync/service"
)

type ReservationHandler struct {
	reservationService service.ReservationService
	validate           *validator.Validate
}
func NewReservationHandler(reservationService service.ReservationService, validate *validator.Validate) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		validate:           validate,
	}
}
func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, _ := c.Get("user_id").(uint)

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", formatValidationErrors(err)))
	}

	reservation, err := h.reservationService.CreateReservation(userID, &req)
	if err != nil {
		if err == repository.ErrZoneFull {
			return c.JSON(http.StatusConflict, dto.ErrorResponse("Parking zone is full, no available spots", nil))
		}
		if err.Error() == "parking zone not found" {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to create reservation", err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse("Reservation confirmed successfully", reservation))
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, _ := c.Get("user_id").(uint)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to retrieve reservations", err.Error()))
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse("My reservations retrieved successfully", reservations))
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, _ := c.Get("user_id").(uint)
	role, _ := c.Get("role").(string)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid reservation ID", nil))
	}

	if err := h.reservationService.CancelReservation(uint(id), userID, role); err != nil {
		switch err.Error() {
		case "reservation not found":
			return c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error(), nil))
		case "forbidden: you can only cancel your own reservations":
			return c.JSON(http.StatusForbidden, dto.ErrorResponse(err.Error(), nil))
		case "only active reservations can be cancelled":
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error(), nil))
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to cancel reservation", err.Error()))
		}
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse("Reservation cancelled successfully", nil))
}


func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to retrieve reservations", err.Error()))
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse("All reservations retrieved successfully", reservations))
}

package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/maftunhasansiam/SpotSync/dto"
	"github.com/maftunhasansiam/SpotSync/service"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validate,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", formatValidationErrors(err)))
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		if err.Error() == "email already registered" {
			return c.JSON(http.StatusConflict, dto.ErrorResponse(err.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Registration failed", err.Error()))
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse("User registered successfully", user))
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
	}
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validation failed", formatValidationErrors(err)))
	}
	result, err := h.authService.Login(&req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse(err.Error(), nil))
	}
	return c.JSON(http.StatusOK, dto.SuccessResponse("Login successful", result))
}

func formatValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errs[e.Field()] = e.Tag()
		}
	}
	return errs
}

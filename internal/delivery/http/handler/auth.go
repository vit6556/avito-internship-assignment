package httphandler

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	tokenTTL    time.Duration
	secure      bool
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService, tokenTTL time.Duration, secure bool) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		tokenTTL:    tokenTTL,
		secure:      secure,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) GetToken(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type must be application/json"})
	}

	var request dto.AuthRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON data"})
	}

	if err := h.validate.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}

	token, err := h.authService.AuthorizeUser(c.Request().Context(), request.Username, request.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid username or password"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to authorize employee"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

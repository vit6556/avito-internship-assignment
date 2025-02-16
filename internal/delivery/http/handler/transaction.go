package httphandler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/vit6556/avito-internship-assignment/internal/delivery/http/dto"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	validate           *validator.Validate
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		validate:           validator.New(),
	}
}

func (h *TransactionHandler) SendCoin(c echo.Context) error {
	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if c.Request().Header.Get("Content-Type") != "application/json" {
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"error": "Content-Type must be application/json"})
	}

	var request dto.SendCoinRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON data"})
	}

	if err := h.validate.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request data"})
	}

	if err := h.transactionService.SendCoins(c.Request().Context(), userID, request.ToUser, request.Amount); err != nil {
		switch err {
		case service.ErrEmployeeNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "employee not found"})
		case service.ErrInsufficientFunds:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "insufficient funds"})
		case service.ErrSelfTransaction:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "sender and recipient cannot be the same user"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send coins"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "coins sent successfully"})
}

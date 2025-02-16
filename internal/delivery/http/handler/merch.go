package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type MerchHandler struct {
	merchService service.MerchService
}

func NewMerchHandler(merchService service.MerchService) *MerchHandler {
	return &MerchHandler{
		merchService: merchService,
	}
}

func (h *MerchHandler) BuyItem(c echo.Context) error {
	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	itemName := c.Param("item")
	if itemName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid merch name"})
	}

	err := h.merchService.BuyItem(c.Request().Context(), userID, itemName)
	if err != nil {
		switch err {
		case service.ErrMerchNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "merch not found"})
		case service.ErrInsufficientFunds:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "insufficient funds"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to buy merch"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "merch purchased successfully"})
}

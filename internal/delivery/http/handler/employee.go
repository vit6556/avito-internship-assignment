package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
}

func NewEmployeeHandler(employeeService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

func (h *EmployeeHandler) GetEmployeeInfo(c echo.Context) error {
	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	employeeInfo, err := h.employeeService.GetEmployeeInfo(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch employee info"})
	}

	return c.JSON(http.StatusOK, employeeInfo)
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

func JWTMiddleware(authService service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := authService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			c.Set("userID", userID)

			return next(c)
		}
	}
}

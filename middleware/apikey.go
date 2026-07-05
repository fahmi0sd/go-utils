package middleware

import (
	"net/http"

	"github.com/fahmi0sd/go-utils/response"
	"github.com/labstack/echo/v4"
)

func APIKeyMiddleware(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if apiKey == "" {
				return next(c)
			}
			if c.Request().Header.Get("X-API-Key") != apiKey {
				return c.JSON(http.StatusUnauthorized, response.Error("unauthorized: invalid api key"))
			}
			return next(c)
		}
	}
}

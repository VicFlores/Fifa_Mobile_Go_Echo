package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func ApiKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if err := godotenv.Load(".env"); err != nil {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		API_KEY := os.Getenv("API_KEY")
		ApiKey := strings.TrimSpace(c.Request().Header.Get("x-api-key"))

		if API_KEY != ApiKey {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalidad credentials"}
		}

		return next(c)
	}
}

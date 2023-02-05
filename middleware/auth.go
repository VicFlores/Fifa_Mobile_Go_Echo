package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/VicFlores/fifa_mobile_API/models"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func ShouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

func CheckAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if !ShouldCheckToken(c.Request().URL.Path) {
			return next(c)
		}

		if err := godotenv.Load(".env"); err != nil {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		JWT_SECRET := os.Getenv("JWT_SECRET")

		tokenString := strings.TrimSpace(c.Request().Header.Get("Authorization"))
		_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET), nil
		})

		if err != nil {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: err.Error()}
		}

		return next(c)
	}
}

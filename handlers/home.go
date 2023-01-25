package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Hello struct {
	Message string `json:"message"`
}

func HomeHandler(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, Hello{
		Message: "Hi Vic Flores 🥵",
	})
}

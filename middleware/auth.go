package middleware

import "github.com/labstack/echo/v4"

type (
	Skipper func(c echo.Context) bool
)

type (
	KeyAuthConfig struct {
		Skipper Skipper

		KeyLookup string

		AuthScheme string

		Validator KeyAuthValidator

		ErrorHandler KeyAuthErrorHandler

		ContinueOnIgnoredError bool
	}

	KeyAuthValidator func(auth string, c echo.Context) (bool, error)

	KeyAuthErrorHandler func(err error, c echo.Context) error
)

var (
	DefaultKeyAuthConfig = KeyAuthConfig{
		Skipper:    DefaultSkipper,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
	}
)

func DefaultSkipper(echo.Context) bool {
	return false
}

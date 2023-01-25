package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	InsertUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	NewUserResponse struct {
		Email string `json:"email"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func SignUpHandler(c echo.Context) (err error) {
	userRequest := InsertUserRequest{}
	echo.New().Validator = &CustomValidator{validator: validator.New()}

	if err = c.Bind(&userRequest); err != nil {
		return err
	}

	if err = c.Validate(&userRequest); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	if userRequest.Email == "" || userRequest.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "missing fields required"}
	}

	/* hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 8)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}

	userModel := models.User{
		Email:    userRequest.Email,
		Password: string(hashedPassword),
	}

	err = repository.SignUp(c.Request().Context(), &userModel)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}

	return c.JSON(http.StatusAccepted, NewUserResponse{
		Email: userModel.Email,
	}) */

	return nil

}

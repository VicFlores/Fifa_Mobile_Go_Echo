package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/VicFlores/fifa_mobile_API/models"
	"github.com/VicFlores/fifa_mobile_API/repository"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	SignUpLoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	NewUserResponse struct {
		Email string `json:"email"`
	}

	SignInResponse struct {
		Token string `json:"token"`
	}
)

var (
	ErrCredentials = echo.NewHTTPError(http.StatusBadRequest, "invalid credentials")
)

func SignUpHandler(c echo.Context) (err error) {
	userRequest := SignUpLoginRequest{}

	if err = c.Bind(&userRequest); err != nil {
		return err
	}

	var validate *validator.Validate = validator.New()

	if err = validate.Struct(&userRequest); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 8)

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
	})
}

func ListUsersHandler(c echo.Context) (err error) {
	var page = uint64(0)
	pageStr := c.Param("page")

	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}
	}

	users, err := repository.ListUsers(c.Request().Context(), page)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, users)
}

func Login(c echo.Context) (err error) {
	loginRequest := SignUpLoginRequest{}
	var validate *validator.Validate = validator.New()

	if err = c.Bind(&loginRequest); err != nil {
		return err
	}

	if err = validate.Struct(&loginRequest); err != nil {
		return &echo.HTTPError{Code: ErrCredentials.Code, Message: ErrCredentials.Message}
	}

	user, err := repository.GetUserByEmail(c.Request().Context(), loginRequest.Email)

	if err != nil {
		return &echo.HTTPError{Code: ErrCredentials.Code, Message: ErrCredentials.Message}
	}

	if user == nil {
		return &echo.HTTPError{Code: ErrCredentials.Code, Message: ErrCredentials.Message}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return &echo.HTTPError{Code: ErrCredentials.Code, Message: ErrCredentials.Message}
	}

	claims := models.AppClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
		},
	}

	if err = godotenv.Load(".env"); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	JWT_SECRET := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return c.JSON(http.StatusOK, SignInResponse{
		Token: tokenString,
	})
}

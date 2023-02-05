package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/VicFlores/fifa_mobile_API/models"
	"github.com/VicFlores/fifa_mobile_API/repository"
	"github.com/VicFlores/fifa_mobile_API/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type (
	UpsertPlayerRequest struct {
		Name     string `json:"name"`
		Position string `json:"position"`
		Club     string `json:"club"`
	}

	NewPlayerResponse struct {
		Name     string `json:"name"`
		Position string `json:"position"`
		Club     string `json:"club"`
	}
)

func InsertPlayerHandler(c echo.Context) (err error) {
	playerRequest := UpsertPlayerRequest{}
	var validate *validator.Validate = validator.New()

	if err = c.Bind(&playerRequest); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	if err = validate.Struct(&playerRequest); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
	}

	playerModel := models.Players{
		Name:     playerRequest.Name,
		Position: playerRequest.Position,
		Club:     playerRequest.Club,
	}

	err = repository.InsertPlayer(c.Request().Context(), &playerModel)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}

	return c.JSON(http.StatusAccepted, NewPlayerResponse{
		Name:     playerModel.Name,
		Position: playerModel.Position,
		Club:     playerModel.Club,
	})
}

func ListPlayersHandler(c echo.Context) (err error) {
	var page = uint64(0)
	pageStr := c.Param("page")

	if pageStr != "" {
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
		}
	}

	players, err := repository.ListPlayers(c.Request().Context(), page)

	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}

	return c.JSON(http.StatusOK, players)
}

func UpdatePlayerHandler(c echo.Context) (err error) {
	tokenString := strings.TrimSpace(c.Request().Header.Get("Authorization"))

	if err = godotenv.Load(".env"); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	JWT_SECRET := os.Getenv("JWT_SECRET")

	env := utils.NewEnvConfig("JWT_SECRET", "PORT")

	if err != nil {
		log.Printf(err.Error())
	}

	log.Println("env", env)

	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: err.Error()}
	}

	if _, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		playerId := c.Param("playerId")
		playerRequest := UpsertPlayerRequest{}

		player, err := repository.GetPlayerById(c.Request().Context(), playerId)

		if err != nil {
			return &echo.HTTPError{Code: ErrCredentials.Code, Message: ErrCredentials.Message}
		}

		if player.Id == 0 && player.Name == "" && player.Club == "" {
			return &echo.HTTPError{Code: ErrCredentials.Code, Message: ErrCredentials.Message}
		}

		if err != nil {
			return &echo.HTTPError{Code: ErrCredentials.Code, Message: ErrCredentials.Message}
		}

		if err = c.Bind(&playerRequest); err != nil {
			return err
		}

		var validate *validator.Validate = validator.New()

		if err = validate.Struct(&playerRequest); err != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}

		playerModel := models.Players{
			Name:     playerRequest.Name,
			Position: playerRequest.Position,
			Club:     playerRequest.Club,
		}

		err = repository.UpdatePlayer(c.Request().Context(), &playerModel, playerId)

		if err != nil {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
		}

		return c.JSON(http.StatusAccepted, NewPlayerResponse{
			Name:     playerModel.Name,
			Position: playerModel.Position,
			Club:     playerModel.Club,
		})
	} else {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}

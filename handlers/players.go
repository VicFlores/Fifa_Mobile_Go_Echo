package handlers

import (
	"net/http"
	"strconv"

	"github.com/VicFlores/fifa_mobile_API/models"
	"github.com/VicFlores/fifa_mobile_API/repository"
	"github.com/labstack/echo/v4"
)

type UpsertPlayerRequest struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Club     string `json:"club"`
}

type NewPlayerResponse struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Club     string `json:"club"`
}

func InsertPlayerHandler(c echo.Context) (err error) {
	playerRequest := UpsertPlayerRequest{}

	if err = c.Bind(&playerRequest); err != nil {
		return err
	}

	if playerRequest.Name == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "missing fields required"}
	}
	if playerRequest.Position == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "missing fields required"}
	}
	if playerRequest.Club == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "missing fields required"}
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

package repository

import (
	"context"

	"github.com/VicFlores/fifa_mobile_API/models"
)

type Repository interface {
	SignUp(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context, page uint64) ([]*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertPlayer(ctx context.Context, player *models.Players) error
	ListPlayers(ctx context.Context, page uint64) ([]*models.Players, error)
	GetPlayerById(ctx context.Context, playerId string) (*models.Players, error)
	UpdatePlayer(ctx context.Context, player *models.Players, playerId string) error
	Close() error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

/* POST */

func InsertPlayer(ctx context.Context, player *models.Players) error {
	return implementation.InsertPlayer(ctx, player)
}

func SignUp(ctx context.Context, user *models.User) error {
	return implementation.SignUp(ctx, user)
}

/* GET */

func ListPlayers(ctx context.Context, page uint64) ([]*models.Players, error) {
	return implementation.ListPlayers(ctx, page)
}

func ListUsers(ctx context.Context, page uint64) ([]*models.User, error) {
	return implementation.ListUsers(ctx, page)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func GetPlayerById(ctx context.Context, playerId string) (*models.Players, error) {
	return implementation.GetPlayerById(ctx, playerId)
}

/* Update */

func UpdatePlayer(ctx context.Context, player *models.Players, playerId string) error {
	return implementation.UpdatePlayer(ctx, player, playerId)
}

func Close() error {
	return implementation.Close()
}

package repository

import (
	"context"

	"github.com/VicFlores/fifa_mobile_API/models"
)

type Repository interface {
	SignUp(ctx context.Context, user *models.User) error
	InsertPlayer(ctx context.Context, player *models.Players) error
	ListPlayers(ctx context.Context, page uint64) ([]*models.Players, error)
	Close() error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

/* CREATE */

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

func Close() error {
	return implementation.Close()
}

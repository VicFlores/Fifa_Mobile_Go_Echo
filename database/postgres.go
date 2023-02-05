package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/VicFlores/fifa_mobile_API/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) SignUp(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO users (email, password) VALUES ($1, $2)",
		user.Email, user.Password)

	return err
}

func (repo *PostgresRepository) InsertPlayer(ctx context.Context, player *models.Players) error {
	_, err := repo.db.ExecContext(ctx,
		"INSERT INTO players (name, position, club) VALUES ($1, $2, $3)",
		player.Name, player.Position, player.Club)

	return err
}

func (repo *PostgresRepository) ListPlayers(ctx context.Context, page uint64) ([]*models.Players, error) {
	rows, err := repo.db.QueryContext(ctx,
		"SELECT id, name, position, club, created_at FROM players LIMIT $1 OFFSET $2", 5, page*5)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var players []*models.Players

	for rows.Next() {
		var player = models.Players{}
		if err = rows.Scan(&player.Id, &player.Name, &player.Position, &player.Club, &player.CreatedAt); err == nil {
			players = append(players, &player)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

func (repo *PostgresRepository) UpdatePlayer(ctx context.Context, player *models.Players, playerId string) error {
	_, err := repo.db.ExecContext(ctx,
		"UPDATE players SET name = $1, position = $2, club = $3 WHERE id = $4",
		player.Name, player.Position, player.Club, playerId)
	return err
}

func (repo *PostgresRepository) ListUsers(ctx context.Context, page uint64) ([]*models.User, error) {
	rows, err := repo.db.QueryContext(ctx,
		"SELECT id, email, created_at FROM users LIMIT $1 OFFSET $2", 5, page*5)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var users []*models.User

	for rows.Next() {
		var user = models.User{}
		if err = rows.Scan(&user.Id, &user.Email, &user.CreatedAt); err == nil {
			users = append(users, &user)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresRepository) GetPlayerById(ctx context.Context, playerId string) (*models.Players, error) {
	rows, err := repo.db.QueryContext(ctx,
		"SELECT id, name, position, club FROM players WHERE id = $1", playerId)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var player = models.Players{}

	for rows.Next() {
		if err = rows.Scan(&player.Id, &player.Name, &player.Position, &player.Club); err == nil {
			return &player, nil
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &player, nil
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

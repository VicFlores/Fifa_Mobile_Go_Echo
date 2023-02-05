package main

import (
	"context"
	"log"
	"os"

	"github.com/VicFlores/fifa_mobile_API/handlers"
	"github.com/VicFlores/fifa_mobile_API/server"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("error loading .env file %v\n", err)
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	server, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatalf("error creating server %v\n", err)
	}

	server.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *echo.Router) {

	r.Add("POST", "/signup", handlers.SignUpHandler)
	r.Add("POST", "/login", handlers.Login)
	r.Add("GET", "/users", handlers.ListUsersHandler)
	r.Add("GET", "/players", handlers.ListPlayersHandler)
	r.Add("POST", "/players", handlers.InsertPlayerHandler)
	r.Add("PUT", "/players/:playerId", handlers.UpdatePlayerHandler)

}

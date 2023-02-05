package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/VicFlores/fifa_mobile_API/database"
	"github.com/VicFlores/fifa_mobile_API/middleware"
	"github.com/VicFlores/fifa_mobile_API/repository"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Config() *Config
}

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Broker struct {
	config *Config
	router *echo.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {

	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	broker := &Broker{
		config: config,
		router: echo.New().Router(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *echo.Router)) {
	newEcho := echo.New()

	newEcho.Use(middleware.ApiKeyMiddleware)
	newEcho.Use(middleware.CheckAuthMiddleware)

	b.router = newEcho.Router()
	binder(b, b.router)

	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)

	if err != nil {
		log.Fatal(err)
	}

	repository.SetRepository(repo)

	s := http.Server{
		Addr:    b.config.Port,
		Handler: newEcho,
	}

	log.Printf("🚀 starting server on port%s 🚀\n", b.config.Port)

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("error starting server:", err)
	} else {
		log.Fatalf("server stopped")
	}
}

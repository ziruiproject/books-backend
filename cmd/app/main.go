package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"starter/config"
	"starter/internal/adapters/api/http/middleware"
	"starter/internal/adapters/database"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to load config file")
	}

	db, err := database.NewPostgresConn()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to connect to database")
	}

	app := &App{}

	router := fiber.New()
	router.Use(middleware.ZerologMiddleware())
	app.Bootstrap(router, db)

	log.Info().Msg("Starting server on :8000")
	for _, route := range router.GetRoutes(true) {
		log.Info().
			Str("method", route.Method).
			Str("path", route.Path).
			//Str("handler", route.Name).
			Msg("Registered Route")
	}
	if err := router.Listen(":8000"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

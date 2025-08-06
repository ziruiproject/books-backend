package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	router.Use(middleware.ZerologMiddleware())
	app.Bootstrap(router, db)

	log.Info().Msgf("Starting server on :%s", config.AppConfig.AppPort)
	for _, route := range router.GetRoutes(true) {
		log.Info().
			Str("method", route.Method).
			Str("path", route.Path).
			Msg("Registered Route")
	}
	if err := router.Listen(fmt.Sprintf(":%s", config.AppConfig.AppPort)); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

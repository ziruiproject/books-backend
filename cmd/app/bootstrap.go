package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/route"
	"starter/internal/adapters/auth"
	"starter/internal/adapters/database"
	"starter/internal/adapters/validator"
	"starter/internal/core/auth"
	"starter/internal/core/user"
	"starter/pkg/hasher"
)

type App struct {
	Handlers   Handlers
	Usecase    Usecase
	Repository Repository
	Route      Route
}

func (app *App) Bootstrap(fiber *fiber.App, db *gorm.DB) {
	app.Repository = *app.NewRepositories()
	app.Usecase = *app.NewUsecases(db)
	app.Handlers = *app.NewHandlers(app.Usecase)
	app.Route = *app.NewRoutes(fiber)
}

type Handlers struct {
	AuthHandler handler.AuthHandler
	UserHandler handler.UserHandler
}

func (app *App) NewHandlers(usecase Usecase) *Handlers {
	return &Handlers{
		AuthHandler: *handler.NewAuthHandler(usecase.AuthUsecase),
		UserHandler: *handler.NewUserHandler(usecase.UserUsecase),
	}
}

type Usecase struct {
	UserUsecase user.Usecase
	AuthUsecase auth.Usecase
}

func (app *App) NewUsecases(db *gorm.DB) *Usecase {
	hasher := hasher.NewBcryptHasher()
	validator := validator.NewValidator()
	tokenGenerator := jwt.NewTokenGenerator()

	userDependency := user.UsecaseDependency{
		DB:             db,
		Hasher:         hasher,
		Validator:      validator,
		UserRepository: app.Repository.UserRepository,
	}
	authDependency := auth.UsecaseDependency{
		DB:             db,
		Hasher:         hasher,
		Validator:      validator,
		UserRepository: app.Repository.UserRepository,
		TokenGenerator: tokenGenerator,
	}
	return &Usecase{
		UserUsecase: user.NewUsecase(userDependency),
		AuthUsecase: auth.NewUsecase(authDependency),
	}
}

type Repository struct {
	UserRepository user.Repository
}

func (app *App) NewRepositories() *Repository {
	return &Repository{
		UserRepository: database.NewUserRepository(),
	}
}

type Route struct {
	UserRoute route.UserRoutes
	AuthRoute route.AuthRoutes
}

func (app *App) NewRoutes(fiber *fiber.App) *Route {
	userRoute := *route.NewUserRoutes(&app.Handlers.UserHandler)
	authRoute := *route.NewAuthRoutes(&app.Handlers.AuthHandler)

	router := fiber.Group("/api/v1")
	userRoute.InstallRoutes(router)
	authRoute.InstallRoutes(router)

	return &Route{
		UserRoute: userRoute,
		AuthRoute: authRoute,
	}
}

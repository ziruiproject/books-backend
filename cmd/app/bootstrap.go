package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"starter/internal/adapters/api/http/handler"
	"starter/internal/adapters/api/http/route"
	"starter/internal/adapters/auth"
	"starter/internal/adapters/database"
	"starter/internal/adapters/storage"
	"starter/internal/adapters/validator"
	"starter/internal/core/auth"
	"starter/internal/core/author"
	"starter/internal/core/book"
	"starter/internal/core/category"
	"starter/internal/core/publisher"
	istorage "starter/internal/core/storage"
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
	AuthHandler      handler.AuthHandler
	UserHandler      handler.UserHandler
	AuthorHandler    handler.AuthorHandler
	CategoryHandler  handler.CategoryHandler
	BookHandler      handler.BookHandler
	PublisherHandler handler.PublisherHandler
	StorageHandler   handler.StorageHandler
}

func (app *App) NewHandlers(usecase Usecase) *Handlers {
	return &Handlers{
		AuthHandler:      *handler.NewAuthHandler(usecase.AuthUsecase),
		UserHandler:      *handler.NewUserHandler(usecase.UserUsecase),
		AuthorHandler:    *handler.NewAuthorHandler(usecase.AuthorUsecase),
		CategoryHandler:  *handler.NewCategoryHandler(usecase.CategoryUsecase),
		PublisherHandler: *handler.NewPublisherHandler(usecase.PublisherUsecase),
		BookHandler:      *handler.NewBookHandler(usecase.BookUsecase),
		StorageHandler:   *handler.NewStorageHandler(usecase.Storage),
	}
}

type Usecase struct {
	UserUsecase      user.Usecase
	AuthUsecase      auth.Usecase
	AuthorUsecase    author.Usecase
	CategoryUsecase  category.Usecase
	PublisherUsecase publisher.Usecase
	BookUsecase      book.Usecase
	Storage          istorage.Storage
}

func (app *App) NewUsecases(db *gorm.DB) *Usecase {
	hasher := hasher.NewBcryptHasher()
	validator := validator.NewValidator()
	storage := storage.NewStorage()
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
	authorDependency := author.UsecaseDependency{
		DB:               db,
		Validator:        validator,
		AuthorRepository: app.Repository.AuthorRepository,
	}
	categoryDependency := category.UsecaseDependency{
		DB:                 db,
		Validator:          validator,
		CategoryRepository: app.Repository.CategoryRepository,
	}
	bookDependency := book.UsecaseDependency{
		DB:             db,
		Validator:      validator,
		Storage:        storage,
		BookRepository: app.Repository.BookRepository,
	}
	publisherDependency := publisher.UsecaseDependency{
		DB:                  db,
		Validator:           validator,
		PublisherRepository: app.Repository.PublisherRepository,
	}

	return &Usecase{
		UserUsecase:      user.NewUsecase(userDependency),
		AuthUsecase:      auth.NewUsecase(authDependency),
		AuthorUsecase:    author.NewUsecase(authorDependency),
		CategoryUsecase:  category.NewUsecase(categoryDependency),
		PublisherUsecase: publisher.NewUsecase(publisherDependency),
		BookUsecase:      book.NewUsecase(bookDependency),
		Storage:          storage,
	}
}

type Repository struct {
	UserRepository      user.Repository
	AuthorRepository    author.Repository
	CategoryRepository  category.Repository
	PublisherRepository publisher.Repository
	BookRepository      book.Repository
}

func (app *App) NewRepositories() *Repository {
	return &Repository{
		UserRepository:      database.NewUserRepository(),
		AuthorRepository:    database.NewAuthorRepository(),
		CategoryRepository:  database.NewCategoryRepository(),
		PublisherRepository: database.NewPublisherRepository(),
		BookRepository:      database.NewBookRepository(),
	}
}

type Route struct {
	UserRoute      route.UserRoutes
	AuthRoute      route.AuthRoutes
	AuthorRoute    route.AuthorRoutes
	CategoryRoute  route.CategoryRoutes
	PublisherRoute route.PublisherRoutes
	BookRoute      route.BookRoutes
	StorageRoute   route.StorageRoutes
}

func (app *App) NewRoutes(fiber *fiber.App) *Route {
	userRoute := *route.NewUserRoutes(&app.Handlers.UserHandler)
	authRoute := *route.NewAuthRoutes(&app.Handlers.AuthHandler)
	authorRoute := *route.NewAuthorRoutes(&app.Handlers.AuthorHandler)
	categoryRoute := *route.NewCategoryRoutes(&app.Handlers.CategoryHandler)
	publisherRoute := *route.NewPublisherRoutes(&app.Handlers.PublisherHandler)
	bookRoute := *route.NewBookRoutes(&app.Handlers.BookHandler)
	storageRoute := *route.NewStorageRoutes(&app.Handlers.StorageHandler)

	router := fiber.Group("/api/v1")
	userRoute.InstallRoutes(router)
	authRoute.InstallRoutes(router)
	authorRoute.InstallRoutes(router)
	categoryRoute.InstallRoutes(router)
	publisherRoute.InstallRoutes(router)
	bookRoute.InstallRoutes(router)
	storageRoute.InstallRoutes(router)

	return &Route{
		UserRoute:      userRoute,
		AuthRoute:      authRoute,
		AuthorRoute:    authorRoute,
		CategoryRoute:  categoryRoute,
		PublisherRoute: publisherRoute,
		BookRoute:      bookRoute,
	}
}

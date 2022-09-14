package api

import (
	"time"

	"github.com/AbdurrahmanA/short-url/api/middleware"
	"github.com/AbdurrahmanA/short-url/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type API struct {
	app             *fiber.App
	port            string
	routes          *routes.Routes
	userHourlyLimit int
	errorHandler    func(c *fiber.Ctx, err error) error
	limiterHandler  func(c *fiber.Ctx) error
}

func NewApi(port string, userHourlyLimit int, routes *routes.Routes, errorHandler func(c *fiber.Ctx, err error) error, limiterHandler func(c *fiber.Ctx) error) *API {
	app := fiber.New(
		fiber.Config{
			ErrorHandler:          errorHandler,
			DisableStartupMessage: true,
		},
	)

	return &API{
		app:             app,
		port:            port,
		routes:          routes,
		userHourlyLimit: userHourlyLimit,
		errorHandler:    errorHandler,
		limiterHandler:  limiterHandler,
	}
}

func (api *API) InitAPI() {

	api.app.Use(cors.New())
	api.app.Use(logger.New())

	api.app.Get("/+", middleware.RedirectShortURLValidation, api.routes.RedirectShortURL)
	api.app.Post("/", limiter.New(limiter.Config{Max: api.userHourlyLimit, LimitReached: api.limiterHandler, Expiration: time.Hour}), middleware.CreateNewShortURLValidation, api.routes.CreateNewShortURL)

	api.app.Listen(":" + api.port)
}

func (api *API) Shutdown() error {
	return api.app.Shutdown()
}

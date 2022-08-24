package api

import (
	"time"

	"github.com/AbdurrahmanA/short-url/api/middleware"
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/AbdurrahmanA/short-url/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusBadRequest

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(response.ErrorResponse{
		Message: err.Error(),
		Status:  code,
	})
}

func limiterHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(response.ErrorResponse{
		Message: "Too many attempts please try again later",
		Status:  fiber.StatusTooManyRequests,
	})
}

func InitAPI(port string, userHourlyLimit int, routes *routes.Routes) {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: errorHandler,
		},
	)

	app.Get("/+", middleware.RedirectShortURLValidation, routes.RedirectShortURL)
	app.Put("/", limiter.New(limiter.Config{Max: userHourlyLimit, LimitReached: limiterHandler, Duration: time.Hour}), middleware.CreateNewShortURLValidation, routes.CreateNewShortURL)

	app.Listen(":" + port)
}

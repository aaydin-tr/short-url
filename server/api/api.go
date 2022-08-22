package api

import (
	"github.com/AbdurrahmanA/short-url/api/middleware"
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/AbdurrahmanA/short-url/api/routes"
	"github.com/gofiber/fiber/v2"
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

func InitAPI(port string, routes *routes.Routes) {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: errorHandler,
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Put("/", middleware.CreateNewShortURLValidation, routes.CreateNewShortURL)

	app.Listen(":" + port)
}

package api

import (
	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
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

func InitAPI(port string) {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: errorHandler,
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Put("/", func(c *fiber.Ctx) error {
		userIP := c.IP()
		body := new(request.NewURL)

		if err := c.BodyParser(body); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		errors := request.ValidateNewURLRequest(*body)
		if errors != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrorResponse{
				Message: "Validation failed",
				Status:  fiber.StatusUnprocessableEntity,
				Errors:  errors,
			})
		}

		return c.SendString(userIP + " " + body.URL)
	})

	app.Listen(":" + port)
}

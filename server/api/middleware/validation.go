package middleware

import (
	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/gofiber/fiber/v2"
)

func CreateNewShortURLValidation(c *fiber.Ctx) error {
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

	c.Locals("body", body)
	return c.Next()
}

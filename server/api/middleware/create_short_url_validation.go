package middleware

import (
	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func validateNewURLRequest(url request.NewURL) []*response.ValidationMessage {
	var errors []*response.ValidationMessage
	err := validate.Struct(url)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element response.ValidationMessage
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			element.Message = helper.MsgForTag(err)
			errors = append(errors, &element)
		}
	}
	return errors
}

func CreateNewShortURLValidation(c *fiber.Ctx) error {
	body := new(request.NewURL)

	if err := c.BodyParser(body); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	errors := validateNewURLRequest(*body)
	if len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrorResponse{
			Message: "Validation failed",
			Status:  fiber.StatusUnprocessableEntity,
			Errors:  errors,
		})
	}

	c.Locals("body", body)
	return c.Next()
}

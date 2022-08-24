package middleware

import (
	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/AbdurrahmanA/short-url/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func validateGetURLRequest(url request.GetURL) []*response.ValidationMessage {
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

func RedirectShortURLValidation(c *fiber.Ctx) error {
	body := request.GetURL{URL: c.Params("+")}
	errors := validateGetURLRequest(body)
	if len(errors) > 0 {
		return c.JSON(response.ErrorResponse{
			Message: "Validation failed",
			Status:  fiber.StatusUnprocessableEntity,
			Errors:  errors,
		})

	}

	c.Locals("shortURL", body.URL)
	return c.Next()
}

package middleware

import (
	"github.com/AbdurrahmanA/short-url/api/request"
	"github.com/AbdurrahmanA/short-url/api/response"
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
			element.Message = err.Error()
			errors = append(errors, &element)
		}
	}
	return errors
}

func validateGetURLRequest(url request.GetURL) []*response.ValidationMessage {
	var errors []*response.ValidationMessage
	err := validate.Struct(url)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element response.ValidationMessage
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			element.Message = msgForTag(err)
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

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "url":
		return "Invalid URL"
	case "alphanum":
		return "Input should be alphanumeric"
	case "len":
		return "Input should be 8 characters long"
	}
	return fe.Error()
}

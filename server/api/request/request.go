package request

import (
	"github.com/AbdurrahmanA/short-url/api/response"
	"github.com/go-playground/validator/v10"
)

type NewURL struct {
	URL string `json:"url" validate:"required,url"`
}

var validate = validator.New()

func ValidateNewURLRequest(url NewURL) []*response.ValidationMessage {
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

package response

import "github.com/AbdurrahmanA/short-url/dto"

type ValidationMessage struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Status  int                  `json:"status"`
	Message string               `json:"message"`
	Errors  []*ValidationMessage `json:"errors"`
}

type SuccessResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    dto.UrlDTO `json:"data"`
}

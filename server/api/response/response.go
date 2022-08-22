package response

import "github.com/AbdurrahmanA/short-url/dto"

type ValidationMessage struct {
	FailedField string
	Tag         string
	Message     string
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

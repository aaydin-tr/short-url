package response

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

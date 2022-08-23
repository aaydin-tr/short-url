package request

type NewURL struct {
	URL string `json:"url" validate:"required,url"`
}

type GetURL struct {
	URL string `validate:"required,len=8,alphanum"`
}

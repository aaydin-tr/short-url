package dto

import (
	"time"

	"github.com/AbdurrahmanA/short-url/model"
)

type UrlDTO struct {
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToUrlDTO(u *model.URL, shortUrlDomain string) UrlDTO {
	return UrlDTO{
		OriginalURL: u.OriginalURL,
		ShortURL:    shortUrlDomain + u.ShortURL,
		CreatedAt:   u.CreatedAt.Time(),
	}
}

package model

type URL struct {
	OwnerIP     string `json:"owner_ip"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

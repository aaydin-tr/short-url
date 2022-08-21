package model

type URL struct {
	OwnerIP     string `bson:"owner_ip"`
	OriginalURL string `bson:"original_url"`
	ShortURL    string `bson:"short_url"`
}

package model

type ShortenResponse struct {
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

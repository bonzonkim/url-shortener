package models

type ShortenRequest struct {
	URL string `json:"url"` // URL to be shortened
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"` // Generated short URL
}

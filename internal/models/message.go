package models

type PunyRequest struct {
	LongURL string `json:"long_url"`
}

type PunyResponse struct {
	ShortURL string `json:"short_url"`
}

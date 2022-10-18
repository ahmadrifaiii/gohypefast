package model

type Payload struct {
	URL string `json:"url"`
}

type Return struct {
	URL      string `json:"url"`
	ShortURL string `json:"shortUrl"`
}

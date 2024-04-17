package shortener

type URL struct {
	URL          string `json:"url"`
	ShortenedURL string `json:"shortened_url"`
}

func NewURL(url string, shortenedURL string) *URL {
	return &URL{URL: url, ShortenedURL: shortenedURL}
}

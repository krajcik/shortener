package shortener

type Url struct {
	Url          string `json:"url"`
	ShortenedUrl string `json:"shortened_url"`
}

func NewUrl(url string, shortenedUrl string) *Url {
	return &Url{Url: url, ShortenedUrl: shortenedUrl}
}

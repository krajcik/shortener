package shortener

type MemRepository struct {
	r map[string]string
}

func NewRepository() *MemRepository {
	return &MemRepository{r: make(map[string]string)}
}

func (r *MemRepository) Save(url *URL) error {
	if _, ok := r.r[url.URL]; ok {
		return ErrAlreadyExists
	}

	r.r[url.URL] = url.ShortenedURL

	return nil
}

func (r *MemRepository) GetByURL(url string) (*URL, error) {
	s, ok := r.r[url]
	if !ok {
		return nil, ErrNotFound
	}

	return NewURL(url, s), nil
}

func (r *MemRepository) GetByShortCode(code string) (*URL, error) {
	for u, s := range r.r {
		if s == code {
			return NewURL(u, s), nil
		}
	}

	return nil, ErrNotFound
}

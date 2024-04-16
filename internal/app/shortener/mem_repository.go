package shortener

type MemRepository struct {
	r map[string]string
}

func NewRepository() *MemRepository {
	return &MemRepository{r: make(map[string]string)}
}

func (r *MemRepository) Save(url *Url) error {
	if _, ok := r.r[url.Url]; ok {
		return AlreadyExistsError
	}

	r.r[url.Url] = url.ShortenedUrl

	return nil
}

func (r *MemRepository) GetByUrl(url string) (*Url, error) {
	s, ok := r.r[url]
	if !ok {
		return nil, NotFoundError
	}

	return NewUrl(url, s), nil
}

func (r *MemRepository) GetByShortCode(code string) (*Url, error) {
	for u, s := range r.r {
		if s == code {
			return NewUrl(u, s), nil
		}
	}

	return nil, NotFoundError
}

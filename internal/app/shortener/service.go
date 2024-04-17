package shortener

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = fmt.Errorf("rep:not found")
	ErrAlreadyExists = fmt.Errorf("rep:already exists")
)

const ShortLen = 11

type Repository interface {
	Save(url *URL) error
	GetByURL(url string) (*URL, error)
	GetByShortCode(code string) (*URL, error)
}

type Service struct {
	r Repository
}

func (s *Service) ShrtByURL(url string) (string, error) {
	shrt, err := s.r.GetByURL(url)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			newURL := NewURL(url, randomString(ShortLen))
			err := s.r.Save(newURL)
			if err != nil {
				return "", err
			}
			return newURL.ShortenedURL, nil
		} else {
			return "", err
		}
	}

	return shrt.ShortenedURL, nil
}

func (s *Service) URLByShrt(shrt string) (string, error) {
	code, err := s.r.GetByShortCode(shrt)
	if err != nil {
		return "", err
	}

	return code.URL, nil
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

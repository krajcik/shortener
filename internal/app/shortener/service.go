package shortener

import (
	"errors"
	"fmt"
)

var (
	NotFoundError      = fmt.Errorf("rep:not found")
	AlreadyExistsError = fmt.Errorf("rep:already exists")
)

const ShortLen = 11

type Repository interface {
	Save(url *Url) error
	GetByUrl(url string) (*Url, error)
	GetByShortCode(code string) (*Url, error)
}

type Service struct {
	r Repository
}

func (s *Service) ShrtByUrl(url string) (string, error) {
	shrt, err := s.r.GetByUrl(url)
	if err != nil {
		if errors.Is(err, NotFoundError) {
			newUrl := NewUrl(url, randomString(ShortLen))
			err := s.r.Save(newUrl)
			if err != nil {
				return "", err
			}
			return newUrl.ShortenedUrl, nil
		} else {
			return "", err
		}
	}

	return shrt.ShortenedUrl, nil
}

func (s *Service) UrlByShrt(shrt string) (string, error) {
	code, err := s.r.GetByShortCode(shrt)
	if err != nil {
		return "", err
	}

	return code.Url, nil
}

func NewService(r Repository) *Service {
	return &Service{r: r}
}

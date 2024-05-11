package shortener

import (
	"context"
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
	SaveBatch(ctx context.Context, urls []*URL) error
	GetByURL(url string) (*URL, error)
	GetByShortCode(code string) (*URL, error)
}

type Service struct {
	r Repository
}

func (s *Service) ShrtBatch(ctx context.Context, urls []string) error {
	var batchToSave []*URL
	for _, url := range urls {
		newURL := NewURL(url, randomString(ShortLen))
		batchToSave = append(batchToSave, newURL)
	}
	err := s.r.SaveBatch(ctx, batchToSave)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ShrtByURL(url string) (string, error) {
	newURL := NewURL(url, randomString(ShortLen))
	err := s.r.Save(newURL)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			byURL, err := s.r.GetByURL(url)
			if err != nil {
				return "", err
			}
			return byURL.ShortenedURL, err
		}
		panic(err)
	}

	return newURL.ShortenedURL, nil
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

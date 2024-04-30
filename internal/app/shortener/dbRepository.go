package shortener

import (
	"database/sql"
	"errors"
)

type DbRepository struct {
	db *sql.DB
}

func NewDbRepository(db *sql.DB) *DbRepository {
	return &DbRepository{db: db}
}

func (d *DbRepository) Save(url *URL) error {
	_, err := d.db.Exec(
		"INSERT INTO url.url(full_url, short_url) VALUES ($1, $2);",
		url.URL, url.ShortenedURL,
	)
	return err
}

func (d *DbRepository) GetByURL(url string) (*URL, error) {
	row := d.db.QueryRow("SELECT full_url, short_url FROM url.url WHERE full_url = $1", url)
	result := &URL{}
	err := row.Scan(&result.URL, &result.ShortenedURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return result, err
}

func (d *DbRepository) GetByShortCode(code string) (*URL, error) {
	row := d.db.QueryRow(`SELECT full_url, short_url FROM url.url WHERE short_url = $1`, code)
	result := &URL{}
	err := row.Scan(&result.URL, &result.ShortenedURL)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	return result, err
}

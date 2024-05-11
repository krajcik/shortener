package shortener

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return fmt.Errorf("%w:%w", ErrAlreadyExists, pgErr)
		}
	}
	return err
}

func (d *DbRepository) SaveBatch(ctx context.Context, urls []*URL) error {
	tx, err := d.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(
		ctx,
		`INSERT INTO url.url(full_url, short_url) VALUES($1, $2) ON CONFLICT (full_url) DO NOTHING;`,
	)
	if err != nil {
		return err
	}
	for _, url := range urls {
		_, err := stmt.ExecContext(
			ctx,
			url.URL,
			url.ShortenedURL,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
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

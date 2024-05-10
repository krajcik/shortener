package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

func Create() (*Params, error) {
	p := &Params{}
	initFlags(p)
	if err := initEnv(p); err != nil {
		return nil, err
	}

	return p, nil
}

func initEnv(p *Params) error {
	err := env.Parse(p)
	if err != nil {
		return err
	}
	return nil
}

func initFlags(p *Params) {
	flag.StringVar(&p.A, "a", ":8080", "The address to listen on for HTTP requests.")
	flag.StringVar(&p.B, "b", "", "Host for shorten url")
	flag.StringVar(&p.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.StringVar(&p.DatabaseDsn, "d", "postgres://postgres:pgpass@localhost:5432/shortener", "Database DSN")
	flag.Parse()
}
